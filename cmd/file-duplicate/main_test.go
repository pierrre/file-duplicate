package main

import (
	"bytes"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestOK(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fl := newFlags()
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := log.New(stderr, "", 0)
	err = run(fl, stdout, l)
	if err != nil {
		t.Fatal(err)
	}
	expectedStdout := filepath.Join(wd, "testdata", "1", "b1") + "\n" + filepath.Join(wd, "testdata", "2", "b2") + "\n\n"
	if stdout.String() != expectedStdout {
		t.Errorf("unexpected stdout: got %q, want %q", stdout.String(), expectedStdout)
	}
	if stderr.String() != "" {
		t.Errorf("unexpected stderr: got %q, want %q", stderr.String(), "")
	}
}

func TestErrorReturn(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fl := newFlags()
	fl.roots = []string{path.Join(wd, "invalid")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := log.New(stderr, "", 0)
	err = run(fl, stdout, l)
	if err == nil {
		t.Fatal("no error")
	}
	if stdout.String() != "" {
		t.Errorf("unexpected stdout: got %q, want %q", stdout.String(), "")
	}
	if stderr.String() != "" {
		t.Errorf("unexpected stderr: got %q, want %q", stderr.String(), "")
	}
}

func TestErrorLog(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fl := newFlags()
	fl.verbose = true
	fl.continueOnError = true
	fl.roots = []string{path.Join(wd, "invalid")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := log.New(stderr, "", 0)
	err = run(fl, stdout, l)
	if err != nil {
		t.Fatal(err)
	}
	if stdout.String() != "" {
		t.Errorf("unexpected stdout: got %q, want %q", stdout.String(), "")
	}
	expectedStderr := "Error: stat " + filepath.Join(wd, "invalid") + string(filepath.Separator) + ".: no such file or directory\n"
	if stderr.String() != expectedStderr {
		t.Errorf("unexpected stderr: got %q, want %q", stderr.String(), expectedStderr)
	}
}
