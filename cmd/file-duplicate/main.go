// Package file-duplicate provides a command line tool to find duplicate files.
package main

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errverbose"
	fileduplicate "github.com/pierrre/file-duplicate"
	"golang.org/x/exp/slog"
)

func main() {
	ctx := context.Background()
	fl := parseFlags()
	l := slog.Default()
	err := run(ctx, fl, os.Stdout, l)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, errverbose.String(err))
		os.Exit(1)
	}
}

func run(ctx context.Context, fl *flags, w io.Writer, l *slog.Logger) error {
	optfs := buildOptions(fl, l)
	err := fileduplicate.Scan(ctx, func(fps []*fileduplicate.File) {
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

func buildOptions(fl *flags, l *slog.Logger) []fileduplicate.Option {
	var optfs []fileduplicate.Option
	fsyss := make([]fs.FS, len(fl.roots))
	for i, root := range fl.roots {
		fsyss[i] = os.DirFS(root)
	}
	optfs = append(optfs, fileduplicate.WithFSs(fsyss))
	if fl.minSize != 0 {
		optfs = append(optfs, fileduplicate.WithMinSize(fl.minSize))
	}
	if fl.continueOnError {
		optfs = append(optfs, fileduplicate.WithErrorHandler(func(ctx context.Context, err error) {
			if fl.verbose {
				l.LogAttrs(ctx, slog.LevelError, errverbose.String(err))
			}
		}))
	}
	return optfs
}
