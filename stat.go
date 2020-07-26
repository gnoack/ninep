package ninep

import "io"

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

func readStat(r io.Reader, s *Stat) error {
	if err := readUint16(r, &s.Type); err != nil {
		return err
	}
	if err := readUint32(r, &s.Dev); err != nil {
		return err
	}
	if err := readUint8(r, &s.QidType); err != nil {
		return err
	}
	if err := readUint32(r, &s.QidVers); err != nil {
		return err
	}
	if err := readUint64(r, &s.QidPath); err != nil {
		return err
	}
	if err := readUint32(r, &s.Mode); err != nil {
		return err
	}
	if err := readUint32(r, &s.Atime); err != nil {
		return err
	}
	if err := readUint32(r, &s.Mtime); err != nil {
		return err
	}
	if err := readUint64(r, &s.Length); err != nil {
		return err
	}
	if err := readString(r, &s.Name); err != nil {
		return err
	}
	if err := readString(r, &s.Uid); err != nil {
		return err
	}
	if err := readString(r, &s.Gid); err != nil {
		return err
	}
	if err := readString(r, &s.Muid); err != nil {
		return err
	}
	return nil
}

func writeStat(w io.Writer, s Stat) error {
	if err := writeUint16(w, s.Type); err != nil {
		return err
	}
	if err := writeUint32(w, s.Dev); err != nil {
		return err
	}
	if err := writeUint8(w, s.QidType); err != nil {
		return err
	}
	if err := writeUint32(w, s.QidVers); err != nil {
		return err
	}
	if err := writeUint64(w, s.QidPath); err != nil {
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
