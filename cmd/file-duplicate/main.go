// Package file-duplicate provides a command line tool to find duplicate files.
package main

import (
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errverbose"
	fileduplicate "github.com/pierrre/file-duplicate"
)

func main() {
	fl := parseFlags()
	l := log.Default()
	err := run(fl, os.Stdout, l)
	if err != nil {
		l.Fatalf("Error: %v", errverbose.Formatter(err))
	}
}

func run(fl *flags, w io.Writer, l *log.Logger) error {
	optfs := buildOptions(fl, l)
	err := fileduplicate.Scan(func(fps []*fileduplicate.File) {
		for _, fp := range fps {
			root := fl.roots[fp.FSIndex]
			p := filepath.Join(root, fp.Path)
			_, _ = io.WriteString(w, p)
			_, _ = io.WriteString(w, "\n")
		}
		_, _ = io.WriteString(w, "\n")
	}, optfs...)
	if err != nil {
		return errors.Wrap(err, "scan")
	}
	return nil
}

func buildOptions(fl *flags, l *log.Logger) []fileduplicate.Option {
	var optfs []fileduplicate.Option
	fsyss := make([]fs.FS, len(fl.roots))
	for i, root := range fl.roots {
		root = filepath.Clean(root)
		if root == "/" {
			root = ""
		}
		fsyss[i] = os.DirFS(root)
	}
	optfs = append(optfs, fileduplicate.WithFSs(fsyss))
	if fl.minSize != 0 {
		optfs = append(optfs, fileduplicate.WithMinSize(fl.minSize))
	}
	if fl.continueOnError {
		optfs = append(optfs, fileduplicate.WithErrorHandler(func(err error) {
			if fl.verbose {
				l.Printf("Error: %v", errverbose.Formatter(err))
			}
		}))
	}
	return optfs
}
