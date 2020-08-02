package ninep

import (
	"errors"
	"io"
)

// TODO: Rename Stat to 'Dir', to be in sync with Plan9 structs.

// Stat represents a directory entry in 9p.
type Stat struct {
	Type uint16 // for kernel use
	Dev  uint32 // for kernel use
	// The type of the file (directory etc.)  represented as a bit
	// vector corresponding to the high 8 bits of the file's mode
	// word.
	// TODO: Use Qid directly here.
	Qid    Qid
	Mode   uint32 // permissions and flags
	Atime  uint32 // last access time
	Mtime  uint32 // last modification time
	Length uint64 // length of file in bytes
	Name   string // file name; must be / if the file is the root
	Uid    string // owner's name
	Gid    string // group's name
	Muid   string // name of the user who last modified the file
}

func readStat(r io.Reader, s *Stat) error {
	var size uint16
	if err := readUint16(r, &size); err != nil {
		return err
	}
	lr := &io.LimitedReader{R: r, N: int64(size)}
	if err := readUint16(lr, &s.Type); err != nil {
		return err
	}
	if err := readUint32(lr, &s.Dev); err != nil {
		return err
	}
	if err := readUint8(lr, &s.Qid.Kind); err != nil {
		return err
	}
	if err := readUint32(lr, &s.Qid.Vers); err != nil {
		return err
	}
	if err := readUint64(lr, &s.Qid.Path); err != nil {
		return err
	}
	if err := readUint32(lr, &s.Mode); err != nil {
		return err
	}
	if err := readUint32(lr, &s.Atime); err != nil {
		return err
	}
	if err := readUint32(lr, &s.Mtime); err != nil {
		return err
	}
	if err := readUint64(lr, &s.Length); err != nil {
		return err
	}
	if err := readString(lr, &s.Name); err != nil {
		return err
	}
	if err := readString(lr, &s.Uid); err != nil {
		return err
	}
	if err := readString(lr, &s.Gid); err != nil {
		return err
	}
	if err := readString(lr, &s.Muid); err != nil {
		return err
	}
	if lr.N > 0 {
		return errors.New("stat is shorter than allocated size")
	}
	return nil
}

func stringSize(s string) uint16 {
	return uint16(2 + len(s))
}

// Size of the given stat struct serialized,
// not including the 2-byte size field.
func statSize(s Stat) (size uint16) {
	size += 39 // fixed part
	size += stringSize(s.Name)
	size += stringSize(s.Uid)
	size += stringSize(s.Gid)
	size += stringSize(s.Muid)
	return size
}

func writeStat(w io.Writer, s Stat) error {
	if err := writeUint16(w, statSize(s)); err != nil {
		return err
	}
	if err := writeUint16(w, s.Type); err != nil {
		return err
	}
	if err := writeUint32(w, s.Dev); err != nil {
		return err
	}
	if err := writeUint8(w, s.Qid.Kind); err != nil {
		return err
	}
	if err := writeUint32(w, s.Qid.Vers); err != nil {
		return err
	}
	if err := writeUint64(w, s.Qid.Path); err != nil {
		return err
	}
	if err := writeUint32(w, s.Mode); err != nil {
		return err
	}
	if err := writeUint32(w, s.Atime); err != nil {
		return err
	}
	if err := writeUint32(w, s.Mtime); err != nil {
		return err
	}
	if err := writeUint64(w, s.Length); err != nil {
		return err
	}
	if err := writeString(w, s.Name); err != nil {
		return err
	}
	if err := writeString(w, s.Uid); err != nil {
		return err
	}
	if err := writeString(w, s.Gid); err != nil {
		return err
	}
	if err := writeString(w, s.Muid); err != nil {
		return err
	}
	return nil
}
