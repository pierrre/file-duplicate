package main

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/pierrre/assert"
)

func TestOK(t *testing.T) {
	ctx := t.Context()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr, nil))
	err = run(ctx, fl, stdout, l)
	assert.NoError(t, err)
	expectedStdout := filepath.Join(wd, "testdata", "1", "b1") + "\n" + filepath.Join(wd, "testdata", "2", "b2") + "\n\n"
	assert.Equal(t, stdout.String(), expectedStdout)
	assert.Zero(t, stderr.String())
}

func TestOKMinSizeZero(t *testing.T) {
	ctx := t.Context()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.minSize = 0
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr, nil))
	err = run(ctx, fl, stdout, l)
	assert.NoError(t, err)
	expectedStdout := filepath.Join(wd, "testdata", "1", "empty1") + "\n" + filepath.Join(wd, "testdata", "2", "empty2") + "\n\n" + filepath.Join(wd, "testdata", "1", "b1") + "\n" + filepath.Join(wd, "testdata", "2", "b2") + "\n\n"
	assert.Equal(t, stdout.String(), expectedStdout)
	assert.Zero(t, stderr.String())
}

func TestOKMinSizePositive(t *testing.T) {
	ctx := t.Context()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.minSize = 3
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr, nil))
	err = run(ctx, fl, stdout, l)
	assert.NoError(t, err)
	assert.Zero(t, stdout.String())
	assert.Zero(t, stderr.String())
}

func TestErrorNoRoots(t *testing.T) {
	ctx := t.Context()
	fl := newFlags()
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr, nil))
	err := run(ctx, fl, stdout, l)
	assert.Error(t, err)
	assert.Zero(t, stdout.String())
	assert.Zero(t, stderr.String())
}

func TestErrorReturn(t *testing.T) {
	ctx := t.Context()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.roots = []string{path.Join(wd, "invalid")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr, nil))
	err = run(ctx, fl, stdout, l)
	assert.Error(t, err)
	assert.Zero(t, stdout.String())
	assert.Zero(t, stderr.String())
}

func TestErrorLog(t *testing.T) {
	ctx := t.Context()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.verbose = true
	fl.continueOnError = true
	fl.roots = []string{path.Join(wd, "invalid")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr, nil))
	err = run(ctx, fl, stdout, l)
	assert.NoError(t, err)
	assert.Zero(t, stdout.String())
	assert.NotZero(t, stderr.String())
}

func TestRunCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr, nil))
	err = run(ctx, fl, stdout, l)
	assert.ErrorIs(t, err, context.Canceled)
	assert.Zero(t, stdout.String())
	assert.Zero(t, stderr.String())
}
