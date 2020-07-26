package ninep

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

const tRead = 123

type callback func()

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

	err error
}

// Peeks at the next available tag without reading it.
func peekTag(r *bufio.Reader) (uint16, error) {
	buf, err := r.Peek(4 + 2 + 2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf[6:8]), nil
}

// Runs the background reader goroutine which dispatches requests.
func (c *clientConn) Run(ctx context.Context) error {
	// TODO: Context cancellation.
	for {
		tag, err := peekTag(c.r)
		if err != nil {
			return fmt.Errorf("Peek error when expecting next message: %w", err)
		}

		c.getReqReader(tag)() // blocking
	}
}

func (c *clientConn) getReqReader(tag uint16) callback {
	c.rrmux.Lock()
	defer c.rrmux.Unlock()

	rr, ok := c.reqReaders[tag]
	if !ok {
		return func() {
			// XXX: Skip next message and log, nothing is registered.
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

func (c *clientConn) Read(fid uint32, offset uint64, buf []byte) (int, error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err := writeTread(c.w, tag.tag, fid, offset, uint32(len(buf)))
	c.wmux.Unlock()

	if err != nil {
		c.err = err
		return 0, err
	}

	tag.await()

	_, data, err := readRread(c.r)
	if err != nil {
		return 0, err
	}

	n := copy(data, buf)
	return n, nil

	// TODO: Would be nice to fill the buf buffer directly instead of copying it over.
	// n := int(r9.Uint16())
	// buf = buf[:n]
	// if _, err := io.ReadFull(r9, buf); err != nil {
	// 	c.err = err
	// 	return 0, c.err
	// }
}

func (c *clientConn) Version(msize uint32, version string) (uint32, string, error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err := writeTversion(c.w, tag.tag, msize, version)
	c.wmux.Unlock()

	if err != nil {
		c.err = err
		return 0, "", err
	}

	tag.await()

	_, rmsize, rversion, err := readRversion(c.r)
	if err != nil {
		c.err = err
		return 0, "", err
	}

	return rmsize, rversion, nil
}

func (c *clientConn) Auth(afid uint32, uname string, aname string) (Qid, error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	err := writeTauth(c.w, tag.tag, afid, uname, aname)
	c.wmux.Unlock()

	if err != nil {
		c.err = err
		return Qid{}, err
	}

	tag.await()

	_, qid, err := readRauth(c.r)
	if err != nil {
		c.err = err
		return qid, err
	}
	return qid, err
}
