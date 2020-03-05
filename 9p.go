package ninep

import (
	"context"
	"io"
	"sync"
)

const tRead = 123

type callback func(uint32, uint16, reader9p)

// p9conn represents a connection to a 9p server.
type p9conn struct {
	tags chan uint16

	wmux sync.Mutex
	w    io.Writer

	r io.Reader

	rrmux      sync.Mutex
	reqReaders map[uint16]callback
}

func (c *p9conn) Run(ctx context.Context) {
	r9 := reader9p{Reader: c.r}
	for {
		size, type9p, tag := r9.Header()

		c.getReqReader(tag)(size, type9p, r9)
	}
}

func (c *p9conn) getReqReader(tag uint16) callback {
	c.rrmux.Lock()
	defer c.rrmux.Unlock()

	rr, ok := c.reqReaders[tag]
	if !ok {
		return func(size uint32, type9p uint16, r9 reader9p) {
			// XXX: Skip next message and log, nothing is registered.
		}
	}

	return rr
}

func (c *p9conn) setReqReader(tag uint16, rr callback) {
	c.rrmux.Lock()
	defer c.rrmux.Unlock()

	c.reqReaders[tag] = rr
}

func (c *p9conn) clearReqReader(tag uint16) {
	c.rrmux.Lock()
	defer c.rrmux.Unlock()

	delete(c.reqReaders, tag)
}

type msgHeader struct {
	size   uint32
	type9p uint16
	r9     reader9p
}

type tagHandle struct {
	tag uint16
	// Reader Run loop sends a msg header for that tag if found.
	ch chan msgHeader
	// The handling function replies back to the reader run loop
	// through this channel.
	done chan struct{}
}

func (h *tagHandle) await() (size uint32, type9p uint16, r9 reader9p) {
	s := <-h.ch
	return s.size, s.type9p, s.r9
}

func (c *p9conn) acquireTag() *tagHandle {
	h := &tagHandle{
		tag:  <-c.tags,
		ch:   make(chan msgHeader),
		done: make(chan struct{}),
	}
	c.setReqReader(h.tag, func(size uint32, type9p uint16, r9 reader9p) {
		// Invoked by reader run loop to read the given message.
		h.ch <- msgHeader{size: size, type9p: type9p, r9: r9}
		<-h.done
	})
	return h
}

func (c *p9conn) releaseTag(h *tagHandle) {
	close(h.done)
	c.clearReqReader(h.tag)
	c.tags <- h.tag
}

func (c *p9conn) Read(fid uint32, offset uint64, buf []byte) (n int, err error) {
	tag := c.acquireTag()
	defer c.releaseTag(tag)

	c.wmux.Lock()
	w := writer9p{Writer: c.w}
	w.Header(4+2+2+4+8+4, tRead, tag.tag)
	w.Uint32(fid)
	w.Uint64(offset)
	w.Uint32(uint32(len(buf)))
	c.wmux.Unlock()

	if w.err != nil {
		// TODO: Error after writing. Handle wedged connection,
		// set p9conn error and return.
	}

	size, type9p, r9 := tag.await()

	_, _ = size, type9p // XXX

	n = int(r9.Uint16())
	buf = buf[:n]
	if _, err := io.ReadFull(r9, buf); err != nil {
		// xxx handle error
		// note reader will be in odd state
		// xxx handle io.ErrUnexpectedEOF.
	}

	return
}
