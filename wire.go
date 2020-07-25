package ninep

//go:generate go run cmd/generate/main.go -o generated.go

import (
	"encoding/binary"
	"io"
)

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
