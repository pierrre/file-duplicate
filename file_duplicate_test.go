package fileduplicate

import (
	"context"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/pierrre/assert"
)

func TestGet(t *testing.T) {
	ctx := t.Context()
	fsys := fstest.MapFS{
		"a": &fstest.MapFile{
			Data: []byte("a"),
		},
		"1/b1": &fstest.MapFile{
			Data: []byte("b"),
		},
		"2/b2": &fstest.MapFile{
			Data: []byte("b"),
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

func TestGetMinSizeZero(t *testing.T) {
	ctx := t.Context()
	fsys := fstest.MapFS{
		"a": &fstest.MapFile{
			Data: []byte("a"),
		},
		"1/b1": &fstest.MapFile{
			Data: []byte("b"),
		},
		"2/b2": &fstest.MapFile{
			Data: []byte("b"),
		},
		"1/empty1": &fstest.MapFile{},
		"2/empty2": &fstest.MapFile{},
	}
	dups, err := Get(ctx, WithFSs([]fs.FS{fsys}), WithMinSize(0))
	assert.NoError(t, err)
	expected := [][]*File{
		{
			{
				FSIndex: 0,
				Path:    "1/empty1",
			},
			{
				FSIndex: 0,
				Path:    "2/empty2",
			},
		},
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

func TestGetCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	fsys := fstest.MapFS{
		"1/b1": &fstest.MapFile{Data: []byte("b")},
		"2/b2": &fstest.MapFile{Data: []byte("b")},
	}
	_, err := Get(ctx, WithFSs([]fs.FS{fsys}))
	assert.ErrorIs(t, err, context.Canceled)
}

func TestScanCanceledDuringCallback(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	fsys := fstest.MapFS{
		"a1": &fstest.MapFile{Data: []byte("a")},
		"a2": &fstest.MapFile{Data: []byte("a")},
		"b1": &fstest.MapFile{Data: []byte("bb")},
		"b2": &fstest.MapFile{Data: []byte("bb")},
	}
	err := Scan(ctx, func(fps []*File) {
		cancel()
	}, WithFSs([]fs.FS{fsys}))
	assert.ErrorIs(t, err, context.Canceled)
}

func TestNewWalkDirFuncCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	opts := newOptions()
	res := make(map[int64][]*File)
	wdf := newWalkDirFunc(ctx, opts, res, 0)
	err := wdf("test", nil, nil)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestGetFilesByHashCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	fsys := fstest.MapFS{
		"1/b1": &fstest.MapFile{Data: []byte("b")},
		"2/b2": &fstest.MapFile{Data: []byte("b")},
	}
	opts := newOptions(WithFSs([]fs.FS{fsys}))
	files := []*File{
		{FSIndex: 0, Path: "1/b1"},
		{FSIndex: 0, Path: "2/b2"},
	}
	_, err := getFilesByHash(ctx, opts, files)
	assert.ErrorIs(t, err, context.Canceled)
}
