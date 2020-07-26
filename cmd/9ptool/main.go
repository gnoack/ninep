package main

import (
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

	components := []string{"plan9"}
	_, err = c.Walk(rootFid, newFid, components)
	if err != nil {
		log.Fatalf("Walk: %v", err)
	}

}
