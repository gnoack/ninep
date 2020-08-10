package ninep

import (
	"encoding/binary"
	"io"
)

func readString(r io.Reader, s *string) error {
	var sz uint16
	if err := binary.Read(r, binary.LittleEndian, &sz); err != nil {
		return err
	}
	buf := make([]byte, sz)
	if err := binary.Read(r, binary.LittleEndian, &buf); err != nil {
		return err
	}
	*s = string(buf)
	return nil
}

func readStringSlice(r io.Reader, ss *[]string) error {
	var size int16
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return err
	}
	*ss = make([]string, 0, size)
	for i := int16(0); i < size; i++ {
		var s string
		if err := readString(r, &s); err != nil {
			return err
		}
		*ss = append(*ss, s)
	}
	return nil
}

func readQIDSlice(r io.Reader, qs *[]QID) error {
	var size uint16
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return err
	}
	*qs = make([]QID, 0, size)
	for i := uint16(0); i < size; i++ {
		var q QID
		if err := readQID(r, &q); err != nil {
			return err
		}
		*qs = append(*qs, q)
	}
	return nil
}

func readQID(r io.Reader, q *QID) error {
	return binary.Read(r, binary.LittleEndian, q)
}

func readByteSlice(r io.Reader, bs *[]byte) error {
	var size uint32 // 4 byte size
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return err
	}
	if uint32(cap(*bs)) < size {
		*bs = make([]byte, size)
	}
	if _, err := io.ReadFull(r, *bs); err != nil {
		return err
	}
	return nil
}

func readUint8(r io.Reader, out *uint8) error {
	return binary.Read(r, binary.LittleEndian, out)
}

func readUint16(r io.Reader, out *uint16) error {
	return binary.Read(r, binary.LittleEndian, out)
}

func readUint32(r io.Reader, out *uint32) error {
	return binary.Read(r, binary.LittleEndian, out)
}

func readUint64(r io.Reader, out *uint64) error {
	return binary.Read(r, binary.LittleEndian, out)
}
