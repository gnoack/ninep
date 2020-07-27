package main

import (
	"fmt"
	"log"

	"github.com/gnoack/ninep"
)

func main() {
	c, err := ninep.Dial("sources")
	if err != nil {
		log.Fatalf("Dial: %v", err)
	}
	_ = c
	//fid := uint32(0)

	rootFid := uint32(0)
	newFid := uint32(1)

	components := []string{"plan9", "NOTICE"}
	qids, err := c.Walk(rootFid, newFid, components)
	if err != nil {
		log.Fatalf("Walk: %v", err)
	}
	fmt.Println("walked to", qids)

	_, _, err = c.Open(newFid, ninep.ORead)
	if err != nil {
		log.Fatalf("Open: %v", err)
	}
	defer c.Clunk(newFid)

	var buf [1000]byte
	n, err := c.Read(newFid, 0, buf[:])
	if err != nil {
		log.Fatalf("Read: %v", err)
	}
	fmt.Println("read", n, "bytes")
	fmt.Println(string(buf[:n]))
}
