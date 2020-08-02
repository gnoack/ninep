package ninep

import "fmt"

// QID in Plan9 is defined in libc.h
type QID struct {
	Path uint64 // uvlong
	Vers uint32 // ulong
	Kind uint8  // uchar
}

func (q QID) String() string { return fmt.Sprintf("{0x%016x %d %d}", q.Path, q.Vers, q.Kind) }
