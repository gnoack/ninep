package ninep

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

type callback func()

// TODO(gnoack): Need a way to close these.
// clientConn represents a connection to a 9p server.
type clientConn struct {
	tags chan uint16

	wmux sync.Mutex
	w    io.Writer

	r *bufio.Reader

	// Callbacks that get called when a message for the given tag
	// is read.
	rrmux      sync.Mutex
	reqReaders map[uint16]callback

	// Connection preferences
	msize uint32
}

// Peeks at the next available tag without reading it.
func peekTag(r *bufio.Reader) (uint16, error) {
	buf, err := r.Peek(4 + 1 + 2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf[5:7]), nil
}

// Runs the background reader goroutine which dispatches requests.
func (c *clientConn) Run(ctx context.Context) error {
	// TODO: Context cancellation.
	for {
		tag, err := peekTag(c.r)
		if err != nil {
			return fmt.Errorf("peek error when expecting next message: %w", err)
		}

		c.getReqReader(tag)() // blocking
	}
}

func (c *clientConn) getReqReader(tag uint16) callback {
	c.rrmux.Lock()
	defer c.rrmux.Unlock()

	rr, ok := c.reqReaders[tag]
	if !ok {
		// Skip message and log, nothing is registered for the tag.
		return func() {
			// TODO: handle errors correctly
			var size uint32
			if err := readUint32(c.r, &size); err != nil {
				return
			}
			buf := make([]byte, size-4)
			n, err := c.r.Read(buf)
			if err != nil || n < int(size-4) {
				return
			}
		}
	}

	return rr
}

func (c *clientConn) setReqReader(tag uint16, rr callback) {
	c.rrmux.Lock()
	defer c.rrmux.Unlock()

	c.reqReaders[tag] = rr
}

func (c *clientConn) clearReqReader(tag uint16) {
	c.rrmux.Lock()
	defer c.rrmux.Unlock()

	delete(c.reqReaders, tag)
}

type tagHandle struct {
	tag uint16
	// Reader Run loop sends a msg header for that tag if found.
	readyToRead chan struct{}
	// The handling function replies back to the reader run loop
	// through this channel.
	doneReading chan struct{}
}

func (h *tagHandle) await() {
	<-h.readyToRead
}

func (c *clientConn) acquireTag() *tagHandle {
	h := &tagHandle{
		tag:         <-c.tags,
		readyToRead: make(chan struct{}),
		doneReading: make(chan struct{}),
	}
	c.setReqReader(h.tag, func() {
		// Invoked by reader run loop to read the given message.
		h.readyToRead <- struct{}{}
		<-h.doneReading
	})
	return h
}

func (c *clientConn) releaseTag(h *tagHandle) {
	close(h.doneReading)
	c.clearReqReader(h.tag)
	c.tags <- h.tag
}

// Read from an open fid.
//
// offset indicates the offset into the file where to read
// buf is the buffer to read into and may not be larger than
// the fid's iounit as returned by Open().
func (c *clientConn) Read(fid uint32, offset uint64, buf []byte) (n uint32, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTread(c.w, tag.tag, fid, offset, uint32(len(buf)))
	c.wmux.Unlock()

	if err != nil {
		return
	}

	tag.await()

	return readRread(c.r, buf)
}

func (c *clientConn) Write(fid uint32, offset uint64, data []byte) (n uint32, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTwrite(c.w, tag.tag, fid, offset, data)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	tag.await()

	return readRwrite(c.r)
}

func (c *clientConn) Version(msize uint32, version string) (rmsize uint32, rversion string, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTversion(c.w, tag.tag, msize, version)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	tag.await()

	return readRversion(c.r)
}

func (c *clientConn) Auth(afid uint32, uname string, aname string) (qid QID, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTauth(c.w, tag.tag, afid, uname, aname)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	tag.await()

	return readRauth(c.r)
}

func (c *clientConn) Attach(fid, afid uint32, uname, aname string) (qid QID, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTattach(c.w, tag.tag, fid, afid, uname, aname)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	tag.await()

	return readRattach(c.r)
}

func (c *clientConn) Walk(fid, newfid uint32, wname []string) (qids []QID, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTwalk(c.w, tag.tag, fid, newfid, wname)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	tag.await()

	return readRwalk(c.r)
}

func (c *clientConn) Stat(fid uint32) (stat Stat, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTstat(c.w, tag.tag, fid)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	tag.await()

	return readRstat(c.r)
}

// The following modes are defined in open(9p) and can be used for
// opening and creating files:
const (
	ORead   = 0x0
	OWrite  = 0x1
	ORdWr   = 0x2
	OExec   = 0x3
	OTrunc  = 0x10 // truncate
	ORClose = 0x40 // delete on clunk
)

func (c *clientConn) Open(fid uint32, mode uint8) (qid QID, iounit uint32, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTopen(c.w, tag.tag, fid, mode)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	tag.await()

	return readRopen(c.r)
}

func (c *clientConn) Clunk(fid uint32) (err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err = writeTclunk(c.w, tag.tag, fid)
	c.wmux.Unlock()

	if err != nil {
		return
	}

	tag.await()

	return readRclunk(c.r)
}
