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
	r9 := reader9p{c.r}
	for {
		size, type9p, tag = r9.Header()

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

func (c *p9conn) Read(fid uint32, offset uint64, buf []byte) (n int, err error) {
	done := make(chan struct{})

	tag := <-c.tags
	defer func() { c.tags <- tag }()

	c.setReqReader(tag, func(size uint32, type9p uint16, r9 reader9p) {
		n := int(r9.Uint16())
		if _, err := io.ReadFull(r9, buf); err != nil {
			// xxx handle error
			// note reader will be in odd state
			// xxx handle io.ErrUnexpectedEOF.
		}

		done <- struct{}{}
	})
	defer func() { c.clearReqReader(tag) }()

	c.wmux.Lock()
	w := writer9p{Writer: c.w}
	w.Header(4+2+2+4+8+4, tRead, tag)
	w.Uint32(fid)
	w.Uint64(offset)
	w.Uint32(uint32(len(buf)))
	c.wmux.Unlock()

	if w.err != nil {
		// TODO: Error after writing. Handle wedged connection,
		// set p9conn error and return.
	}

	<-done
	return
}
