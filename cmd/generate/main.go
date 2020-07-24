package main

import (
	"bufio"
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

// returns type, variable name
func getInfo(s string) (string, string) {
	parts := strings.SplitN(s, "[", 2)
	name := parts[0]
	switch {
	case strings.HasSuffix(s, "[1]"):
		return "uint8", name
	case strings.HasSuffix(s, "[2]"):
		return "uint16", name
	case strings.HasSuffix(s, "[4]"):
		return "uint32", name
	case strings.HasSuffix(s, "[8]"):
		return "uint64", name
	case strings.HasSuffix(s, "[13]"):
		return "Qid", name
	case strings.HasSuffix(s, "[s]"):
		return "string", name
	case strings.HasPrefix(s, "T") || strings.HasPrefix(s, "R"):
		return "uint16", "msgType"
	case s == "stat[n]":
		return "Stat", name
	case strings.HasSuffix(s, "[count[4]]"):
		return "[]byte", name
	case s == "nwname*(wname[s])":
		return "[]string", "nwnames"
	case s == "nwqid*(qid[13])":
		return "[]Qid", "qids"
	default:
		log.Fatalf("unknown type: %q", s)
	}
	return "", ""
}

func getType(s string) string {
	t, _ := getInfo(s)
	return t
}

func printWriteFunc(ss []string) {
	name := ss[1]
	var msgType string
	fmt.Print("func Write" + name + "(w io.Writer, ")
	for i, s := range ss {
		t, n := getInfo(s)
		if n == "msgType" {
			msgType = s
			continue
		}
		fmt.Print(n + " " + t)
		if i < len(ss)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println(") error {")
	for _, s := range ss {
		t, n := getInfo(s)
		funcname := fmt.Sprintf("Write%v", strings.Title(t))
		if t == "[]string" {
			funcname = "WriteStringSlice"
		}
		if t == "[]Qid" {
			funcname = "WriteStringSlice"
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

func main() {
	out := make(chan []string)
	go collect(out)

	for ss := range out {
		ss = conflate(ss)
		fmt.Println()
		printComment(ss)
		printWriteFunc(ss)
	}
}
