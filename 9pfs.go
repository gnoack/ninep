package ninep

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type file struct {
	FID    uint32
	cc     *clientConn
	offset uint64
	iounit uint32
	QID    QID
}

func (f *file) Read(p []byte) (n int, err error) {
	// Truncate read to iounit size if necessary.
	if uint32(len(p)) > f.iounit {
		p = p[:f.iounit]
	}
	count, err := f.cc.Read(context.TODO(), f.FID, f.offset, p)
	if err != nil {
		return 0, err
	}
	if count == 0 && len(p) > 0 {
		return 0, io.EOF
	}
	f.offset += uint64(count) // XXX: Check overflow?
	return int(count), nil
}

func (f *file) Stat() (info os.FileInfo, err error) {
	stat, err := f.cc.Stat(context.TODO(), f.FID)
	return &statFileInfo{s: stat}, err
}

func (f *file) ReadDir(n int) (infos []os.FileInfo, err error) {
	if !f.QID.IsDirectory() {
		return nil, errors.New("not a directory")
	}
	br := bufio.NewReader(f)
	unlimited := n <= 0
	for i := 0; i < n || unlimited; i++ {
		var stat Stat
		if err := readStat(br, &stat); err != nil {
			if unlimited && err == io.EOF {
				err = nil
			}
			return infos, err
		}
		infos = append(infos, &statFileInfo{s: stat})
	}
	return infos, nil
}

func (f *file) Close() error {
	return f.cc.Clunk(context.TODO(), f.FID)
}

// TODO: Double check that the mode bits match.
type statFileInfo struct{ s Stat }

func (fi *statFileInfo) Name() string       { return fi.s.Name }
func (fi *statFileInfo) Size() int64        { return int64(fi.s.Length) }
func (fi *statFileInfo) Mode() os.FileMode  { return os.FileMode(fi.s.Mode) }
func (fi *statFileInfo) ModTime() time.Time { return time.Unix(int64(fi.s.Mtime), 0) }
func (fi *statFileInfo) IsDir() bool        { return (fi.s.Mode & ModeDir) != 0 }
func (fi *statFileInfo) Sys() interface{}   { return fi.s }

type FS struct {
	cc      *clientConn
	nextFID uint32
	rootFID uint32
}

// Open opens a file for reading.
func (f *FS) Open(name string) (*file, error) {
	// TODO: Verify name format.
	components := strings.Split(name, "/")
	if len(name) == 0 {
		components = nil
	}

	// TODO: Track used FIDs instead of just cycling.
	f.nextFID++
	_, err := f.cc.Walk(context.TODO(), f.rootFID, f.nextFID, components)
	if err != nil {
		return nil, fmt.Errorf("9p walk: %w", err)
	}

	qid, iounit, err := f.cc.Open(context.TODO(), f.nextFID, ORead)
	if err != nil {
		return nil, fmt.Errorf("9p open: %w", err)
	}
	// If iounit is 0, we need to fall back to connection message
	// size - 24.
	if iounit == 0 {
		iounit = f.cc.msize - 24
	}

	return &file{FID: f.nextFID, cc: f.cc, iounit: iounit, QID: qid}, nil
}
