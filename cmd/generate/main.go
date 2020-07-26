package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func collect(out chan []string) {
	walkFunc := func(path string, info os.FileInfo, err error) error {
		_ = info
		if !strings.HasSuffix(path, ".9p") {
			return nil
		}

		base := strings.TrimSuffix(filepath.Base(path), ".9p")
		if base == "0intro" {
			return nil
		}

		r, err := os.Open(path)
		must(err)
		defer r.Close()

		s := bufio.NewScanner(r)

		// Skip until past the SYNOPSIS header
		for s.Scan() && strings.TrimSpace(s.Text()) != ".SH SYNOPSIS" {
		}

		var parts []string
		for s.Scan() && !strings.HasPrefix(strings.TrimSpace(s.Text()), ".SH ") {
			l := strings.TrimSpace(s.Text())
			if strings.HasPrefix(l, ".ta ") || l == ".br" || l == ".PP" {
				if len(parts) > 0 {
					out <- parts
					parts = nil
				}
				continue // new function command
			}
			if strings.HasPrefix(l, ".IR ") {
				l = strings.TrimPrefix(l, ".IR ")
				l = strings.ReplaceAll(l, " ", "")
				parts = append(parts, l)
				continue
			}
			if strings.HasPrefix(l, ".B ") {
				l = strings.TrimPrefix(l, ".B ")
				l = strings.ReplaceAll(l, " ", "")
				parts = append(parts, l)
				continue
			}
			must(fmt.Errorf("unrecognized line: %q", l))
		}
		must(s.Err())
		if len(parts) > 0 {
			out <- parts
		}
		return nil
	}
	defer close(out)
	filepath.Walk("/usr/lib/plan9/man/man9", walkFunc)
}

func printComment(ss []string) {
	fmt.Print("//")
	for _, s := range ss {
		fmt.Print(" " + s)
	}
	fmt.Println()
}

// returns type, variable name, size calculation code
func getInfo(s string) (string, string, string) {
	parts := strings.SplitN(s, "[", 2)
	name := parts[0]
	switch {
	case strings.HasSuffix(s, "[1]"):
		return "uint8", name, "1"
	case strings.HasSuffix(s, "[2]"):
		return "uint16", name, "2"
	case strings.HasSuffix(s, "[4]"):
		return "uint32", name, "4"
	case strings.HasSuffix(s, "[8]"):
		return "uint64", name, "8"
	case strings.HasSuffix(s, "[13]"):
		return "Qid", name, "13"
	case strings.HasSuffix(s, "[s]"):
		return "string", name, fmt.Sprintf("(2 + len(%v))", name)
	case strings.HasPrefix(s, "T") || strings.HasPrefix(s, "R"):
		return "uint16", "msgType", "2"
	case s == "stat[n]":
		return "Stat", name, fmt.Sprintf("(39 + 8 + len(%v.Name) + len(%v.Uid) + len(%v.Gid) + len(%v.Muid))", name, name, name, name)
	case strings.HasSuffix(s, "[count[4]]"):
		return "[]byte", name, fmt.Sprintf("(4 + len(%v))", name)
	case s == "nwname*(wname[s])":
		return "[]string", "nwnames", "stringSliceSize(nwnames)"
	case s == "nwqid*(qid[13])":
		return "[]Qid", "qids", "(2 + 13*len(qids))"
	default:
		log.Fatalf("unknown type: %q", s)
	}
	return "", "", ""
}

func dontReturnTag(name string) bool {
	// We don't want to return the tag when reading reply data;
	// the tags are already peeked in advance of reading by the 9p
	// protocol layer.
	return name[0] == 'R'
}

func printReadFunc(ss []string) {
	name := ss[1]
	fmt.Print("func read" + name + "(r io.Reader) (")
	for _, s := range ss {
		t, n, _ := getInfo(s)
		if n == "msgType" || n == "size" {
			continue
		}
		if n == "tag" && dontReturnTag(name) {
			continue
		}
		fmt.Print(n + " " + t + ", ")
	}
	fmt.Println("err error) {")
	fmt.Println("  var size uint32")
	for _, s := range ss {
		t, n, _ := getInfo(s)
		funcname := fmt.Sprintf("read%v", strings.Title(t))
		if t == "[]string" {
			funcname = "readStringSlice"
		}
		if t == "[]Qid" {
			funcname = "readQidSlice"
		}
		if t == "[]byte" {
			funcname = "readByteSlice"
		}
		if n == "msgType" {
			fmt.Println("  var msgType uint16")
		}
		if n == "tag" && dontReturnTag(name) {
			fmt.Println("  var tag uint16")
		}
		fmt.Printf("  if err = %v(r, &%v); err != nil {\n", funcname, n)
		fmt.Println("    return")
		fmt.Println("  }")
		if n == "msgType" {
			if name[0] == 'R' {
				fmt.Println("  if msgType == Rerror {")
				fmt.Println("    // XXX Read error contents")
				fmt.Println("    err = backendError")
				fmt.Println("    return")
				fmt.Println("  }")
			}
			fmt.Println("  if msgType !=", name, "{")
			fmt.Println("    err = unexpectedMsgError")
			fmt.Println("    return")
			fmt.Println("  }")
		}
	}
	fmt.Println("  return")
	fmt.Println("}")
}

func printWriteFunc(ss []string) {
	name := ss[1]
	var msgType string
	fmt.Print("func write" + name + "(w io.Writer, ")
	for i, s := range ss {
		t, n, _ := getInfo(s)
		// msgType is fixed for each method
		if n == "msgType" {
			msgType = s
			continue
		}
		// size is calculated dynamically based on other parameters
		if n == "size" {
			continue
		}
		fmt.Print(n + " " + t)
		if i < len(ss)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println(") error {")
	// Size calculation
	fmt.Print("  size := uint32(")
	for i, s := range ss {
		_, _, sz := getInfo(s)
		fmt.Print(sz)
		if i < len(ss)-1 {
			fmt.Print(" + ")
		}
	}
	fmt.Println(")")

	for _, s := range ss {
		t, n, _ := getInfo(s)
		funcname := fmt.Sprintf("write%v", strings.Title(t))
		if t == "[]string" {
			funcname = "writeStringSlice"
		}
		if t == "[]Qid" {
			funcname = "writeQidSlice"
		}
		if t == "[]byte" {
			funcname = "writeByteSlice"
		}
		if n == "msgType" {
			n = msgType // resolve to constant directly
		}
		fmt.Printf("  if err := %v(w, %v); err != nil {\n", funcname, n)
		fmt.Println("    return err")
		fmt.Println("  }")
	}
	fmt.Println("  return nil")
	fmt.Println("}")
}

// Conflate:
// count[4] data[count] => data[count[4]]
func conflate(in []string) []string {
	var o []string
	for i := 0; i < len(in); i++ {
		// byte buffers
		if in[i] == "count[4]" {
			if i+1 < len(in) && in[i+1] == "data[count]" {
				o = append(o, "data[count[4]]")
				i++
				continue
			}
		}
		// directory names for walk
		if in[i] == "nwname[2]" {
			if i+1 < len(in) && in[i+1] == "nwname*(wname[s])" {
				o = append(o, "nwname*(wname[s])")
				i++
				continue
			}
		}
		// qids for walk response
		if in[i] == "nwqid[2]" {
			if i+1 < len(in) && in[i+1] == "nwqid*(qid[13])" {
				o = append(o, "nwqid*(qid[13])")
				i++
				continue
			}
		}
		o = append(o, in[i])
	}
	return o
}

var outfile = flag.String("o", "/dev/stdout", "output file")

func main() {
	flag.Parse()
	f, err := os.Create(*outfile)
	must(err)
	defer f.Close()
	os.Stdout = f

	out := make(chan []string)
	go collect(out)

	fmt.Println(`package ninep

import "io"`)

	for ss := range out {
		ss = conflate(ss)
		fmt.Println()
		printComment(ss)
		printWriteFunc(ss)

		fmt.Println()
		printComment(ss)
		printReadFunc(ss)
	}
}
