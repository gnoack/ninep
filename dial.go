package ninep

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// nofid is the fid value used to indicate absence of a FID,
// e.g. to pass as afid when no authentication is required.
const nofid uint32 = ^uint32(0)

// notag is the tag value used in absence of a tag,
// e.g. during authentication
const notag uint16 = ^uint16(0)

func dialNet(service string) (net.Conn, error) {
	if service == "sources" {
		return net.Dial("tcp", "sources.9p.io:564")
	}
	if strings.HasPrefix(service, "localhost:") {
		return net.Dial("tcp", service)
	}
	sessionDir := fmt.Sprintf("ns.%s.%s", os.Getenv("USER"), os.Getenv("DISPLAY"))
	return net.Dial("unix", filepath.Join("/tmp", sessionDir, service))
}

func versionRPC(c net.Conn, wantVersion string, wantMsize uint32) (msize uint32, vErr error) {
	if err := writeTversion(c, notag, wantMsize, wantVersion); err != nil {
		return 0, err
	}
	msize, version, err := readRversion(c)
	if err != nil {
		return 0, fmt.Errorf("version(%q, %q): %w", wantMsize, wantVersion, err)
	}

	if wantMsize < msize {
		return 0, fmt.Errorf("server wanted too high msize of %v", msize)
	}

	if version != wantVersion {
		return 0, fmt.Errorf("mismatching version: %q != %q", version, wantVersion)
	}
	return msize, nil
}

type DialFSOpts struct {
	DialOpts
	AttachOpts
}

// DialFS dials a 9p client connection and directly attaches to it.
func DialFS(service string, opts DialFSOpts) (dFS *FS, dErr error) {
	cc, err := Dial(service, opts.DialOpts)
	if err != nil {
		return nil, err
	}

	return Attach(cc, opts.AttachOpts)
}

type DialOpts struct {
	Concurrency uint16
}

// Dial establishes a 9p client connection and returns it.
func Dial(service string, opts DialOpts) (dConn *ClientConn, dErr error) {
	if opts.Concurrency == 0 {
		opts.Concurrency = 256
	}

	// Dial.
	netConn, err := dialNet(service)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dConn != nil {
			return
		}
		netConn.Close()
	}()

	// Check version and negotiate msize.
	msize, err := versionRPC(netConn, "9P2000", 8192)
	if err != nil {
		return nil, err
	}

	// Build client connection.
	ctx, cancelCause := context.WithCancelCause(context.Background())
	cc := &ClientConn{
		tags:       make(chan uint16, opts.Concurrency),
		conn:       netConn,
		reqReaders: make(map[uint16]callback),
		msize:      msize,
		cancel:     cancelCause,
	}
	// Fill tag queue.
	for i := uint16(0); i < opts.Concurrency; i++ {
		cc.tags <- i
	}

	cc.wg.Add(1)
	go func() {
		defer cc.wg.Done()
		err := cc.run(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return // OK
			}
			// TODO: How to report error correctly?
			log.Fatalf("9p client: run(): %v", err)
		}
	}()

	return cc, nil
}

type Authenticator func(io.ReadWriter) error

type AttachOpts struct {
	// Username to attach with
	Uname string

	// The remote file system to attach to
	Aname string

	// Authenticator
	Authenticator Authenticator
}

// Attach opens a file system from an already-open client connection.
func Attach(cc *ClientConn, opts AttachOpts) (fsys *FS, err error) {
	// Attempt auth
	afid := cc.fidPool.Acquire()

	qid, err := cc.Auth(context.Background(), afid, opts.Uname, opts.Aname)
	switch {
	case err != nil && strings.HasSuffix(err.Error(), "authentication not required"):
		// Authentication not required.
		cc.fidPool.Release(afid)
		afid = nofid
	case err != nil:
		cc.fidPool.Release(afid)
		return nil, err
	case opts.Authenticator == nil:
		cc.fidPool.Release(afid)
		return nil, errors.New("no means to authenticate")
	default:
		// TODO: This should not duplicate the code from OpenFile().
		iounit := cc.msize - 24
		authfile := &file{FID: afid, cc: cc, iounit: iounit, QID: qid}
		defer authfile.Close()

		err := opts.Authenticator(authfile)
		if err != nil {
			return nil, err
		}
	}

	// Attach.
	fid := cc.fidPool.Acquire()
	defer func() {
		if fsys != nil {
			return
		}
		cc.fidPool.Release(fid)
	}()
	_, err = cc.Attach(context.Background(), fid, afid, opts.Uname, opts.Aname)
	if err != nil {
		return nil, err
	}

	return &FS{cc: cc, rootFID: fid}, nil
}
