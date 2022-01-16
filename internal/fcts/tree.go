package fcts

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

func cleanPath(file string) string {
	return strings.TrimSuffix(file, ".gpg")
}

func Tree(w io.Writer, dir string) error {
	tree, err := newTree(dir)
	for _, line := range tree {
		// FIXME: we need a nice formatting for the tree output
		fmt.Fprintln(w, line)
	}
	return err
}

func TreeFind(w io.Writer, dir string, terms []string) error {
	tree, err := newTree(dir)
	if err != nil {
		return err
	}
	for _, line := range tree {
		var found bool = false
		for _, term := range terms {
			if strings.Contains(strings.Join(line, ""), term) {
				found = true
			}
		}
		if found {
			fmt.Fprintln(w, line)
		}
	}
	return err
}

func newTree(dir string) ([][]string, error) {
	var tree [][]string
	buildTree := func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			dir, file := filepath.Split(strings.TrimPrefix(path, dir))
			var line []string
			for _, dir := range filepath.SplitList(dir) {
				line = append(line, filepath.Base(dir))
			}
			line = append(line, cleanPath(file))
			tree = append(tree, line)
		}
		return nil
	}
	err := filepath.Walk(dir, buildTree)
	return tree, err
}
