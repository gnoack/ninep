package ninep

import (
	"context"
	"fmt"
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

type dialOptions struct {
	concurrency uint16
}

type dialOpt func(*dialOptions)

func WithConcurrency(concurrency uint16) dialOpt {
	return func(c *dialOptions) {
		c.concurrency = concurrency
	}
}

func Dial(service string, opts ...dialOpt) (dFS *FS, dErr error) {
	cc, err := dial9pConn(service, opts...)
	if err != nil {
		return nil, err
	}

	// Attach.
	var (
		afid  = nofid
		uname = "user"
		aname = ""
	)
	fid := cc.fidPool.Acquire()
	defer func() {
		if dFS != nil {
			return
		}
		cc.fidPool.Release(fid)
	}()
	_, err = cc.Attach(context.Background(), fid, afid, uname, aname)
	if err != nil {
		return nil, err
	}

	return &FS{cc: cc, rootFID: fid}, nil
}

// dial9pConn establishes a 9p client connection and returns it.
func dial9pConn(service string, opts ...dialOpt) (dConn *ClientConn, dErr error) {
	options := dialOptions{
		concurrency: 256,
	}
	for _, opt := range opts {
		opt(&options)
	}

	// Dial
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
		tags:       make(chan uint16, options.concurrency),
		conn:       netConn,
		reqReaders: make(map[uint16]callback),
		msize:      msize,
		cancel:     cancelCause,
	}
	// Fill tag queue.
	for i := uint16(0); i < options.concurrency; i++ {
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
