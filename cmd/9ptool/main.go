package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gnoack/ninep"
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage\n")
	fmt.Fprintf(flag.CommandLine.Output(), "     %s CMD PATH\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "e.g. %s cat sources/plan9/NOTICE\n\n", os.Args[0])
	flag.PrintDefaults()
}

func parsePositionalArgs() (cmd string, service string, path string) {
	if len(flag.Args()) != 2 {
		log.Fatal("Could not parse positional args; want [cmd] [path]")
	}
	cmd = flag.Args()[0]
	arg := flag.Args()[1]
	parts := strings.SplitN(arg, "/", 2)
	service = parts[0]
	if len(parts) == 2 {
		path = parts[1]
	}
	return
}

func formatStat(stat os.FileInfo) string {
	sys := stat.Sys().(ninep.Stat)
	return fmt.Sprintf("%s %8d %8s %8s %s",
		stat.Mode().String(), stat.Size(), sys.UID, sys.GID, stat.Name())
}

func main() {
	// For better RPC latency debugging, log microseconds.
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	flag.Usage = usage
	flag.Parse()
	cmd, service, path := parsePositionalArgs()

	c, err := ninep.DialFS(service)
	if err != nil {
		log.Fatalf("DialFS(%q): %v", service, err)
	}

	r, err := c.Open(path)
	if err != nil {
		log.Fatalf("Open: %v", err)
	}
	defer r.Close()

	switch cmd {
	case "cat":
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatalf("Read: %v", err)
		}
		os.Stdout.Write(buf)

	case "stat":
		stat, err := r.Stat()
		if err != nil {
			log.Fatalf("Stat: %v", err)
		}
		fmt.Println(formatStat(stat))

	case "ls":
		infos, err := r.ReadDir(0)
		if err != nil {
			log.Fatalf("ReadDir: %v", err)
		}
		for _, info := range infos {
			fmt.Println(formatStat(info))
		}
	}
}
