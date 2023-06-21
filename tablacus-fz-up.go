package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		cur   string
		root  string
		filer string
	)
	flag.StringVar(&cur, "cur", "", "current dir path")
	flag.StringVar(&root, "root", "", "root path")
	flag.StringVar(&filer, "filer", "explorer.exe", "filer")
	flag.Parse()
	os.Exit(run(cur, root, filer))
}

func run(c string, root string, filer string) int {
	found := getDirs(c, root)
	if len(found) < 1 {
		return 0
	}
	idx, err := fuzzyfinder.Find(found, func(i int) string {
		return formatRel(root, found[i])
	})
	if err != nil {
		if err == fuzzyfinder.ErrAbort {
			return 0
		}
		return 1
	}
	src := found[idx]
	if fi, err := os.Stat(src); err == nil && fi.IsDir() {
		exec.Command(filer, src).Start()
	} else {
		return 1
	}
	return 0
}

func fromPath(s string) []string {
	return strings.Split(s, string(os.PathSeparator))
}
func toPath(ss []string) string {
	return strings.Join(ss, string(os.PathSeparator))
}

func getDirs(cur string, root string) []string {
	elems := fromPath(cur)
	var found []string
	for i := 1; i < len(elems); i++ {
		if i == len(elems)-1 {
			return []string{}
		}
		p := toPath(elems[0 : len(elems)-i])
		if p == root {
			break
		}
		found = append(found, p)
	}
	return found
}

func formatRel(root string, p string) string {
	h := "~"
	rel, err := filepath.Rel(root, p)
	if err != nil {
		return h
	}
	elems := fromPath(rel)
	if 1 < len(elems) {
		return h + string(os.PathSeparator) + toPath(elems[1:])
	}
	return h
}
