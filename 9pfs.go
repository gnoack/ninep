package ninep

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type file struct {
	fid    uint32
	cc     *clientConn
	offset uint64
	iounit uint32
}

func (f *file) Read(p []byte) (n int, err error) {
	// Truncate read to iounit size if necessary.
	if uint32(len(p)) > f.iounit {
		p = p[:f.iounit]
	}
	count, err := f.cc.Read(f.fid, f.offset, p)
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
	stat, err := f.cc.Stat(f.fid)
	return &statFileInfo{s: stat}, err
}

func (f *file) ReadDir(n int) (infos []os.FileInfo, err error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}
	if !stat.IsDir() {
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
	return f.cc.Clunk(f.fid)
}

// TODO: Double check that the mode bits match.
type statFileInfo struct{ s Stat }

func (fi *statFileInfo) Name() string       { return fi.s.Name }
func (fi *statFileInfo) Size() int64        { return int64(fi.s.Length) }
func (fi *statFileInfo) Mode() os.FileMode  { return os.FileMode(fi.s.Mode) }
func (fi *statFileInfo) ModTime() time.Time { return time.Unix(int64(fi.s.Mtime), 0) }
func (fi *statFileInfo) IsDir() bool        { return (fi.s.Mode & ModeDir) != 0 }
func (fi *statFileInfo) Sys() interface{}   { return fi.s }

type fs struct {
	cc      *clientConn
	nextFid uint32
	rootFid uint32
}

func (f *fs) Open(name string) (*file, error) {
	// TODO: Verify name format.
	// TODO: Track used FIDs instead of just cycling.
	f.nextFid++

	components := strings.Split(name, "/")
	_, err := f.cc.Walk(f.rootFid, f.nextFid, components)
	if err != nil {
		return nil, fmt.Errorf("9p walk: %w", err)
	}

	_, iounit, err := f.cc.Open(f.nextFid, ORead)
	if err != nil {
		return nil, fmt.Errorf("9p open: %w", err)
	}
	// TODO: If iounit is 0, do we need to fall back to
	// connection message size - 24?
	if iounit == 0 {
		return nil, fmt.Errorf("9p open: iounit is 0")
	}

	return &file{fid: f.nextFid, cc: f.cc, iounit: iounit}, nil
}
