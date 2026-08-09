// Package file-duplicate provides a command-line tool to find duplicate files.
package main

import (
	"context"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errverbose"
	fileduplicate "github.com/pierrre/file-duplicate"
	"github.com/pierrre/go-libs/unsafeio"
)

func main() {
	os.Exit(mainRun())
}

func mainRun() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	fl := parseFlags()
	l := slog.Default()
	err := run(ctx, fl, os.Stdout, l)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			l.LogAttrs(ctx, slog.LevelError, errverbose.String(err))
		}
		return 1
	}
	return 0
}

func run(ctx context.Context, fl *flags, w io.Writer, l *slog.Logger) error {
	optfs := buildOptions(fl, l)
	err := fileduplicate.Scan(ctx, func(fps []*fileduplicate.File) {
		for _, fp := range fps {
			root := fl.roots[fp.FSIndex]
			p := filepath.Join(root, fp.Path)
			_, _ = unsafeio.WriteString(w, p)
			_, _ = unsafeio.WriteString(w, "\n")
		}
		_, _ = unsafeio.WriteString(w, "\n")
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
	optfs = append(optfs, fileduplicate.WithFSs(fsyss), fileduplicate.WithMinSize(fl.minSize))
	if fl.continueOnError {
		optfs = append(optfs, fileduplicate.WithErrorHandler(func(ctx context.Context, err error) {
			if fl.verbose {
				l.LogAttrs(ctx, slog.LevelError, errverbose.String(err))
			}
		}))
	}
	return optfs
}
