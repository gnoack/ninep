package ninep

import "fmt"

// in Plan9, defined in libc.h
type Qid struct {
	Path uint64 // uvlong
	Vers uint32 // ulong
	Kind uint8  // uchar
}

func (q Qid) String() string { return fmt.Sprintf("{0x%016x %d %d}", q.Path, q.Vers, q.Kind) }

// TODO: Rename this to 'Dir', to be in sync with Plan9 structs.
type Stat struct {
	Type uint16 // for kernel use
	Dev  uint32 // for kernel use
	// The type of the file (directory etc.)  represented as a bit
	// vector corresponding to the high 8 bits of the file's mode
	// word.
	// TODO: Use Qid directly here.
	QidType uint8
	QidVers uint32 // version number for given path
	QidPath uint64 // the file server's unique ID for the file
	Mode    uint32 // permissions and flags
	Atime   uint32 // last access time
	Mtime   uint32 // last modification time
	Length  uint64 // length of file in bytes
	Name    string // file name; must be / if the file is the root
	Uid     string // owner's name
	Gid     string // group's name
	Muid    string // name of the user who last modified the file
}
