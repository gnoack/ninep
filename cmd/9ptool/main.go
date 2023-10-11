package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gnoack/ninep"
)

var (
	uname = flag.String("uname", os.Getenv("USER"), "Username to try to attach with")
	aname = flag.String("aname", "", "File system to attach to (may be empty)")
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
	service, path, _ = strings.Cut(arg, "/")
	return
}

func formatStat(stat fs.FileInfo) string {
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

	fsys, err := ninep.DialFS(service, ninep.DialFSOpts{
		AttachOpts: ninep.AttachOpts{
			Uname:         *uname,
			Aname:         *aname,
			Authenticator: interact,
		},
	})
	if err != nil {
		log.Fatalf("DialFS(%q): %v", service, err)
	}
	defer fsys.Close()

	switch cmd {
	case "cat":
		r, err := fsys.Open(path)
		if err != nil {
			log.Fatalf("Open: %v", err)
		}
		defer r.Close()
		buf, err := io.ReadAll(r)
		if err != nil {
			log.Fatalf("Read: %v", err)
		}
		os.Stdout.Write(buf)

	case "write":
		f, err := fsys.OpenFile(path, ninep.OWrite)
		if err != nil {
			log.Fatalf("OpenFile: %v", err)
		}
		defer f.Close()

		w, ok := f.(io.Writer)
		if !ok {
			log.Fatalf("OpenFile did not return writeable file: %v", f)
		}
		n, err := io.Copy(w, os.Stdin)
		if err != nil {
			log.Fatalf("io.Copy: %v", err)
		}
		fmt.Printf("%d bytes written.\n", n)

	case "rpc":
		f, err := fsys.OpenFile(path, ninep.ORdWr)
		if err != nil {
			log.Fatalf("OpenFile: %v", err)
		}
		defer f.Close()

		rw, ok := f.(io.ReadWriter)
		if !ok {
			log.Fatalf("not a read-writeable file: %v", f)
		}
		interact(rw)

	case "stat":
		stat, err := fs.Stat(fsys, path)
		if err != nil {
			log.Fatalf("Stat: %v", err)
		}
		fmt.Println(formatStat(stat))

	case "ls":
		entries, err := fs.ReadDir(fsys, path)
		if err != nil {
			log.Fatalf("ReadDir: %v", err)
		}
		for _, e := range entries {
			info, err := e.Info()
			if err != nil {
				log.Fatalf("Info(): %v", err)
			}
			fmt.Println(formatStat(info))
		}
	}
}

// TODO: Clean up this mess - this is probably wrong.
func interact(f io.ReadWriter) error {
	time.Sleep(500 * time.Millisecond)
	stdinR := bufio.NewReader(os.Stdin)
	var buf [8192]byte
	for {
		for {
			n, err := f.Read(buf[:])
			if err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Println("err:", err)
					time.Sleep(500 * time.Millisecond)
					break
				}
				return err
			}
			if n == 0 {
				break
			}
			fmt.Println("*** " + string(buf[:n]))
		}
		line, err := stdinR.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		line = strings.TrimSuffix(line, "\n") + "\x00"
		_, err = f.Write([]byte(line))
		if err != nil {
			if strings.HasPrefix(err.Error(), "bad rpc verb") {
				log.Printf("Error: %v", err)
				continue
			}
			return err
		}
	}
}
