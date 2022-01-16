package fcts

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

func cleanPath(file string) string {
	return strings.TrimSpace(strings.TrimSuffix(file, ".gpg"))
}

func Tree(w io.Writer, dir string) error {
	tree, err := newTree(dir)
	// FIXME: we need a nice formatting for the tree output
	printTree(w, tree, "")
	return err
}

func TreeFind(w io.Writer, dir string, terms []string) error {
	tree, err := newTree(dir)
	if err != nil {
		return err
	}
	var results [][]string
	for _, line := range tree {
		var found bool = false
		for _, term := range terms {
			if strings.Contains(strings.Join(line, ""), term) {
				found = true
			}
		}
		if found {
			results = append(results, line)
		}
	}
	// printTree(w, results, []bool{})
	printTree(w, results, "")
	return err
}

func newTree(dir string) ([][]string, error) {
	var tree [][]string
	buildTree := func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			dir, file := filepath.Split(strings.TrimPrefix(path, dir))
			var line []string
			line = append(line, splitDir(dir)...)

			if !strings.HasPrefix(cleanPath(file), ".") {
				line = append(line, cleanPath(file))
				tree = append(tree, line)
			}
		}
		return nil
	}
	err := filepath.Walk(dir, buildTree)
	return tree, err
}

func splitDir(dir string) []string {
	// FIXME: this is platform specific
	dirs := strings.Split(dir, "/")
	var results []string
	for _, dir := range dirs {
		if dir != "" {
			results = append(results, strings.TrimSpace(dir))
		}
	}
	return results
}

func printTree(w io.Writer, tree [][]string, pre string) error {
	content := make(map[string][][]string)
	for _, line := range tree {
		if len(line) == 0 {
			continue
		}
		if node, ok := content[line[0]]; ok && len(line) > 0 {
			node = append(node, line[1:])
			content[line[0]] = node
		} else if len(line) > 0 {
			content[line[0]] = [][]string{line[1:]}
		}
	}
	var keys []string
	for key, _ := range content {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for i, key := range keys {
		nextPre := pre
		dirFmt := pre
		if i == (len(keys) - 1) {
			nextPre += "  "
			dirFmt += "\u2514\u2500"
		} else {
			nextPre += "\u2502 "
			dirFmt += "\u251C\u2500"
		}
		fmt.Fprintf(w, "%s%s%s\n", dirFmt, "", key)
		printTree(w, content[key], nextPre)
	}
	return nil
}
