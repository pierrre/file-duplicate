package fileduplicate

import (
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestGet(t *testing.T) {
	fsys := fstest.MapFS{
		"a": &fstest.MapFile{
			Data: []byte(("a")),
		},
		"1/b1": &fstest.MapFile{
			Data: []byte(("b")),
		},
		"2/b2": &fstest.MapFile{
			Data: []byte(("b")),
		},
		"1/empty1": &fstest.MapFile{},
		"2/empty2": &fstest.MapFile{},
	}
	dups, err := Get(WithFSs([]fs.FS{fsys}))
	if err != nil {
		t.Fatal(err)
	}
	if len(dups) != 1 {
		t.Fatalf("unexpected elements: got %d, want %d", len(dups), 1)
	}
	if len(dups[0]) != 2 {
		t.Fatalf("unexpected duplicates: got %d, want %d", len(dups[0]), 2)
	}
	if dups[0][0].FSIndex != 0 {
		t.Fatalf("unexpected fs index: got %d, want %d", dups[0][0].FSIndex, 0)
	}
	if dups[0][0].Path != "1/b1" {
		t.Fatalf("unexpected path: got %s, want %s", dups[0][0].Path, "1/b1")
	}
	if dups[0][1].FSIndex != 0 {
		t.Fatalf("unexpected fs index: got %d, want %d", dups[0][1].FSIndex, 0)
	}
	if dups[0][1].Path != "2/b2" {
		t.Fatalf("unexpected path: got %s, want %s", dups[0][1].Path, "2/b2")
	}
}
