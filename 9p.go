package ninep

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"
)

var errConnShutdown = errors.New("connection shutdown")

type msgHeader struct {
	size    uint32
	msgType uint8
	tag     uint16
}

func (h *msgHeader) serialize() (out [7]byte) {
	enc := binary.LittleEndian
	enc.PutUint32(out[0:4], h.size)
	out[4] = byte(h.msgType)
	enc.PutUint16(out[5:7], h.tag)
	return out
}

func (h *msgHeader) readerFrom(r io.Reader) io.Reader {
	hdrBuf := h.serialize()
	hdrReader := bytes.NewBuffer(hdrBuf[:])
	bodyReader := io.LimitReader(r, int64(h.size-7))
	return io.MultiReader(hdrReader, bodyReader)
}

type callback func(d msgHeader)

// ClientConn represents a connection to a 9p server.
type ClientConn struct {
	tags chan uint16

	wmux sync.Mutex // Write mutex.
	conn io.ReadWriteCloser

	// Callbacks that get called when a message for the given tag is read.
	rrmux      sync.Mutex // Mutex for reqReaders.
	reqReaders map[uint16]callback

	// Connection preferences
	msize uint32

	// Shutdown helpers
	cancel func(error)
	wg     sync.WaitGroup

	// Thread-safe pool of FIDs to use
	fidPool fidPool
}

func readHeader(r io.Reader) (hdr msgHeader, err error) {
	var buf [7]byte
	_, err = io.ReadFull(r, buf[:])
	if err != nil {
		return
	}
	return msgHeader{
		size:    binary.LittleEndian.Uint32(buf[0:4]),
		msgType: uint8(buf[4]),
		tag:     binary.LittleEndian.Uint16(buf[5:7]),
	}, nil
}

// run runs the background reader goroutine which dispatches requests.
func (c *ClientConn) run(ctx context.Context) error {
	// TODO: The context cancelation here is poor. We should be
	// able to return immediately, not only after reading the next
	// message header.
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		hdr, err := readHeader(c.conn)
		// TODO: Use limit reader
		if err != nil {
			if context.Cause(ctx) == errConnShutdown {
				return nil
			}
			return fmt.Errorf("peek error when expecting next message: %w", err)
		}

		c.getReqReader(hdr.tag)(hdr) // blocking
	}
}

func (c *ClientConn) getReqReader(tag uint16) callback {
	c.rrmux.Lock()
	defer c.rrmux.Unlock()

	rr, ok := c.reqReaders[tag]
	if !ok {
		// Skip message and log, nothing is registered for the tag.
		return func(hdr msgHeader) {
			// TODO: handle errors correctly
			err := skip(c.conn, int(hdr.size-7))
			if err != nil {
				return
			}
		}
	}

	return rr
}

func (c *ClientConn) setReqReader(tag uint16, rr callback) {
	c.rrmux.Lock()
	defer c.rrmux.Unlock()

	c.reqReaders[tag] = rr
}

func (c *ClientConn) clearReqReader(tag uint16) {
	c.rrmux.Lock()
	defer c.rrmux.Unlock()

	delete(c.reqReaders, tag)
}

// Close closes the 9p connection.
func (c *ClientConn) Close() error {
	if c.cancel == nil {
		return nil
	}

	defer c.wg.Wait()
	c.cancel(errConnShutdown)
	c.cancel = nil
	return c.conn.Close()
}

type tagHandle struct {
	tag uint16
	// Reader run loop sends a msg header for that tag if found.
	readyToRead chan msgHeader
	// The handling function replies back to the reader run loop
	// through this channel.
	doneReading chan struct{}
	// Parent ClientConn
	conn *ClientConn
}

func (h *tagHandle) awaitHdr(ctx context.Context) (msgHeader, error) {
	select {
	case hdr := <-h.readyToRead:
		return hdr, nil
	case <-ctx.Done():
		return msgHeader{}, ctx.Err()
	}
}

// Await the response for the given tag. On success, returns a reader
// for the response message (bounded to size). Returns ctx.Err() on
// early cancelation.
func (h *tagHandle) await(ctx context.Context) (io.Reader, error) {
	hdr, err := h.awaitHdr(ctx)
	if err != nil {
		return nil, err
	}
	return hdr.readerFrom(h.conn.conn), nil
}

func (c *ClientConn) acquireTag() *tagHandle {
	h := &tagHandle{
		conn:        c,
		tag:         <-c.tags,
		readyToRead: make(chan msgHeader),
		doneReading: make(chan struct{}),
	}
	c.setReqReader(h.tag, func(hdr msgHeader) {
		// Invoked by reader run loop to read the given message.
		h.readyToRead <- hdr
		<-h.doneReading
	})
	return h
}

func (c *ClientConn) releaseTag(h *tagHandle) {
	close(h.doneReading)
	c.clearReqReader(h.tag)
	c.tags <- h.tag
}

// Read from an open fid.
//
// offset indicates the offset into the file where to read.
// buf is the buffer to read into and may not be larger than
// the fid's iounit as returned by Open().
func (c *ClientConn) Read(ctx context.Context, fid uint32, offset uint64, buf []byte) (n uint32, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTread(c.conn, tag.tag, fid, offset, uint32(len(buf)))
	c.wmux.Unlock()

	if err != nil {
		return
	}

	r, err := tag.await(ctx)
	if err != nil {
		c.Flush(tag.tag)
		return
	}

	return readRread(r, buf)
}

func (c *ClientConn) Write(ctx context.Context, fid uint32, offset uint64, data []byte) (n uint32, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTwrite(c.conn, tag.tag, fid, offset, data)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	r, err := tag.await(ctx)
	if err != nil {
		c.Flush(tag.tag)
		return
	}

	return readRwrite(r)
}

func (c *ClientConn) Walk(ctx context.Context, fid, newfid uint32, wname []string) (qids []QID, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTwalk(c.conn, tag.tag, fid, newfid, wname)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	r, err := tag.await(ctx)
	if err != nil {
		c.Flush(tag.tag)
		return
	}

	return readRwalk(r)
}

func (c *ClientConn) Stat(ctx context.Context, fid uint32) (stat Stat, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTstat(c.conn, tag.tag, fid)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	r, err := tag.await(ctx)
	if err != nil {
		c.Flush(tag.tag)
		return
	}

	return readRstat(r)
}

// Modes for opening and creating files, as defined in open(9p).
const (
	ORead   = 0x0
	OWrite  = 0x1
	ORdWr   = 0x2
	OExec   = 0x3
	OTrunc  = 0x10 // truncate
	ORClose = 0x40 // delete on clunk
)

func (c *ClientConn) Open(ctx context.Context, fid uint32, mode uint8) (qid QID, iounit uint32, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTopen(c.conn, tag.tag, fid, mode)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	r, err := tag.await(ctx)
	if err != nil {
		c.Flush(tag.tag)
		return
	}

	return readRopen(r)
}

func (c *ClientConn) Clunk(ctx context.Context, fid uint32) (err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTclunk(c.conn, tag.tag, fid)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	r, err := tag.await(ctx)
	if err != nil {
		c.Flush(tag.tag)
		return
	}

	return readRclunk(r)
}

// TODO: Do callers need to check the error?
func (c *ClientConn) Flush(oldtag uint16) (err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTflush(c.conn, tag.tag, oldtag)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	// Note: This may not error. Servers must repond to flush.
	r, _ := tag.await(context.Background())

	return readRflush(r)
}

func (c *ClientConn) Attach(ctx context.Context, fid uint32, afid uint32, uname string, aname string) (qid QID, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTattach(c.conn, tag.tag, fid, afid, uname, aname)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	r, err := tag.await(context.Background())
	if err != nil {
		c.Flush(tag.tag)
		return
	}

	return readRattach(r)
}

func (c *ClientConn) Auth(ctx context.Context, afid uint32, uname, aname string) (qid QID, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTauth(c.conn, tag.tag, afid, uname, aname)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	r, err := tag.await(context.Background())
	if err != nil {
		c.Flush(tag.tag)
		return
	}

	return readRauth(r)
}
