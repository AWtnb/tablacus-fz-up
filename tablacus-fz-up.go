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
	b := getBase(c, root)
	idx, err := fuzzyfinder.Find(found, func(i int) string {
		rel, _ := filepath.Rel(b, found[i])
		return rel
	})
	if err != nil {
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

func getBase(c string, root string) string {
	elems := fromPath(c)
	rd := fromPath(root)
	if len(elems) < len(rd) {
		return root
	}
	return toPath(elems[0 : len(rd)+1])
}
