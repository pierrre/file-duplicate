package fileduplicate

import (
	"context"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/pierrre/assert"
)

func TestGet(t *testing.T) {
	ctx := context.Background()
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
	dups, err := Get(ctx, WithFSs([]fs.FS{fsys}))
	assert.NoError(t, err)
	expected := [][]*File{
		{
			{
				FSIndex: 0,
				Path:    "1/b1",
			},
			{
				FSIndex: 0,
				Path:    "2/b2",
			},
		},
	}
	assert.DeepEqual(t, dups, expected)
}
