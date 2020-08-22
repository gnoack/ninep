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

// Nofid is the fid value used to indicate absence of a FID,
// e.g. to pass as afid when no authentication is required.
const Nofid uint32 = ^uint32(0)

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

func handshake(c net.Conn) (msize uint32, err error) {
	uname, aname := "user", ""
	wantVersion := "9P2000"
	var wantMsize uint32 = 8192
	rootFID := uint32(0) // TODO: Dynamically acquire FIDs somehow

	if err := writeTversion(c, 0xffff, wantMsize, wantVersion); err != nil {
		return 0, err
	}
	msize, version, err := readRversion(c)
	if err != nil {
		return 0, fmt.Errorf("version(%q, %q): %w", wantMsize, wantVersion, err)
	}

	if msize < wantMsize {
		// TODO: Fall back to server-provided msize if needed
		return 0, fmt.Errorf("server wanted too high msize of %v", msize)
	}
	if version != wantVersion {
		return 0, fmt.Errorf("mismatching version: %q != %q", version, wantVersion)
	}

	// Afid is nofid when the client doesn't want to authenticate.
	afid := Nofid

	// XXX: Authentication step

	if err := writeTattach(c, 1, rootFID, afid, uname, aname); err != nil {
		return 0, err
	}
	_, err = readRattach(c)
	if err != nil {
		return 0, fmt.Errorf("attach(): %w", err)
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

func DialFS(service string, opts ...dialOpt) (*fs, error) {
	cc, err := Dial(service, opts...)
	if err != nil {
		return nil, err
	}
	return &fs{cc: cc}, nil
}

func Dial(service string, opts ...dialOpt) (*clientConn, error) {
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

	// Handshake
	msize, err := handshake(netConn)
	if err != nil {
		netConn.Close()
		return nil, err
	}

	// Build client connection.
	cc := &clientConn{
		tags:       make(chan uint16, options.concurrency),
		w:          netConn,
		r:          netConn,
		reqReaders: make(map[uint16]callback),
		msize:      msize,
	}
	go func() {
		for i := uint16(0); i < options.concurrency; i++ {
			cc.tags <- i
		}
	}()

	go func() {
		err := cc.Run(context.Background()) // TODO: Cancellation
		if err != nil {
			// TODO: How to report error correctly?
			log.Fatalf("9p client: Run(): %v", err)
		}
	}()

	return cc, nil
}
