package ninep

import (
	"encoding/binary"
	"io"
)

// writer9p is a wrapper around io.Writer with 9p specific helpers.
type writer9p struct {
	io.Writer
	err error
}

func (w *writer9p) applyErr(err error) {
	if err != nil && w.err == nil {
		w.err = err
	}
}

func (w *writer9p) Header(size uint32, type9p uint16, tag uint16) {
	w.Uint32(size)
	w.Uint16(type9p)
	w.Uint16(tag)
}

func (w *writer9p) Uint64(v uint64) {
	w.applyErr(binary.Write(*w, binary.LittleEndian, v))
}

func (w *writer9p) Uint32(v uint32) {
	w.applyErr(binary.Write(*w, binary.LittleEndian, v))
}

func (w *writer9p) Uint16(v uint16) {
	w.applyErr(binary.Write(*w, binary.LittleEndian, v))
}

// reader9p is a wrapper around io.Reader with 9p specific helpers.
type reader9p struct {
	io.Reader
	err error
}

func (r *reader9p) Header() (size uint32, type9p uint16, tag uint16) {
	size = r.Uint32()
	type9p = r.Uint16()
	tag = r.Uint16()
	return
}

func (r *reader9p) applyErr(err error) {
	if err != nil && r.err == nil {
		r.err = err
	}
}

func (r *reader9p) Uint32() (v uint32) {
	r.applyErr(binary.Read(*r, binary.LittleEndian, &v))
	return
}

func (r *reader9p) Uint16() (v uint16) {
	r.applyErr(binary.Read(*r, binary.LittleEndian, &v))
	return
}

func (r *reader9p) String() (s string) {
	var n uint16
	r.applyErr(binary.Read(r, binary.LittleEndian, &n))
	buf := make([]byte, n)
	r.applyErr(binary.Read(r, binary.LittleEndian, &buf))
	return string(buf)
}
