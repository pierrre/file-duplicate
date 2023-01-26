package fileduplicate

import (
	"io/fs"
	"testing"
	"testing/fstest"

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
