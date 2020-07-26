package ninep

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// The Nofid constant is used to indicate absence of a FID,
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

func Dial(service string) (*clientConn, error) {
	// Dial options
	tagCapacity := uint16(20) // TODO
	uname, aname := "user", ""
	wantVersion := "9P2000"
	var wantMsize uint32 = 4000
	rootFid := uint32(0) // TODO: Dynamically acquire FIDs somehow

	// Dial
	netConn, err := dialNet(service)
	if err != nil {
		return nil, err
	}

	cc := &clientConn{
		tags:       make(chan uint16, tagCapacity),
		w:          netConn,
		r:          bufio.NewReader(netConn),
		reqReaders: make(map[uint16]callback),
	}
	go func() {
		for i := uint16(0); i < tagCapacity; i++ {
			cc.tags <- i
		}
	}()

	go cc.Run(context.Background()) // TODO: Cancellation

	msize, version, err := cc.Version(wantMsize, wantVersion)
	if err != nil {
		return nil, fmt.Errorf("Version(%q, %q): %w", wantMsize, wantVersion, err)
	}
	if msize < wantMsize {
		// TODO: Fall back to server-provided msize if needed
		return nil, fmt.Errorf("Server wanted too high msize of %v", msize)
	}
	if version != wantVersion {
		return nil, fmt.Errorf("Mismatching version: %q != %q", version, wantVersion)
	}

	// Afid is nofid when the client doesn't want to authenticate.
	afid := Nofid

	// XXX: Authentication step

	_, err = cc.Attach(rootFid, afid, uname, aname)
	if err != nil {
		return nil, fmt.Errorf("Attach(): %w", err)
	}
	return cc, nil
}
