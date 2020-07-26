package ninep

import "fmt"

// in Plan9, defined in libc.h
type Qid struct {
	Path uint64 // uvlong
	Vers uint32 // ulong
	Kind uint8  // uchar
}

func (q Qid) String() string { return fmt.Sprintf("{0x%016x %d %d}", q.Path, q.Vers, q.Kind) }
