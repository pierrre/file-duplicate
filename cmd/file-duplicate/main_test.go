package main

import (
	"bytes"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/assert/ext/davecghspew"
	"github.com/pierrre/assert/ext/pierrrecompare"
	"github.com/pierrre/assert/ext/pierrreerrors"
)

func init() {
	pierrrecompare.Configure()
	davecghspew.ConfigureDefault()
	pierrreerrors.Configure()
}

func TestOK(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := log.New(stderr, "", 0)
	err = run(fl, stdout, l)
	assert.NoError(t, err)
	expectedStdout := filepath.Join(wd, "testdata", "1", "b1") + "\n" + filepath.Join(wd, "testdata", "2", "b2") + "\n\n"
	assert.Equal(t, stdout.String(), expectedStdout)
	assert.StringEmpty(t, stderr.String())
}

func TestErrorReturn(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.roots = []string{path.Join(wd, "invalid")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := log.New(stderr, "", 0)
	err = run(fl, stdout, l)
	assert.Error(t, err)
	assert.StringEmpty(t, stdout.String())
	assert.StringEmpty(t, stderr.String())
}

func TestErrorLog(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.verbose = true
	fl.continueOnError = true
	fl.roots = []string{path.Join(wd, "invalid")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := log.New(stderr, "", 0)
	err = run(fl, stdout, l)
	assert.NoError(t, err)
	assert.StringEmpty(t, stdout.String())
	assert.StringNotEmpty(t, stderr.String())
}
