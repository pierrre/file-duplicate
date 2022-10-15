// Package fileduplicate provides utilities to find duplicate files.
package fileduplicate

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"sort"

	"github.com/pierrre/errors"
)

type options struct {
	fss          []fs.FS
	minSize      int64
	errorHandler func(error)
}

func newOptions(optfs ...Option) *options {
	opts := &options{
		minSize: 1,
	}
	for _, optf := range optfs {
		optf(opts)
	}
	return opts
}

// Option represents an option.
type Option func(*options)

// WithFSs is an option that defines the filesystems to scan.
func WithFSs(fsyss []fs.FS) Option {
	return func(o *options) {
		o.fss = fsyss
	}
}

// WithMinSize is an option that defines the minimum file size to consider.
func WithMinSize(minSize int64) Option {
	return func(o *options) {
		o.minSize = minSize
	}
}

// WithErrorHandler is an option that defines the error handler.
//
// If it is defined, the error handler is called for each error, otherwise the error is returned.
func WithErrorHandler(f func(error)) Option {
	return func(o *options) {
		o.errorHandler = f
	}
}

// File represents a file.
type File struct {
	// FSIndex is the index of the filesystem where the file is located.
	FSIndex int
	// Path is the path of the file in the filesystem.
	Path string
}

// Get returns the duplicated files.
func Get(optfs ...Option) ([][]*File, error) {
	var res [][]*File
	err := Scan(func(fps []*File) {
		res = append(res, fps)
	}, optfs...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Scan scans for duplicated files.
//
// The onDuplicates function is called for each duplicated file.
func Scan(onDuplicates func([]*File), optfs ...Option) error {
	opts := newOptions(optfs...)
	sameSizeFiles, err := getSameSizeFiles(opts)
	if err != nil {
		return err
	}
	for _, fpsS := range sameSizeFiles {
		filesByHash, err := getFilesByHash(opts, fpsS)
		if err != nil {
			return err
		}
		for _, fpsH := range filesByHash {
			if len(fpsH) <= 1 {
				continue
			}
			onDuplicates(fpsH)
		}
	}
	return nil
}

func getSameSizeFiles(opts *options) ([][]*File, error) {
	filesBySize, err := getFilesBySize(opts)
	if err != nil {
		return nil, err
	}
	var sizes []int64
	for size, fps := range filesBySize {
		if len(fps) > 1 {
			sizes = append(sizes, size)
		}
	}
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] < sizes[j]
	})
	res := make([][]*File, 0, len(sizes))
	for _, size := range sizes {
		fps := filesBySize[size]
		res = append(res, fps)
	}
	return res, nil
}

func getFilesBySize(opts *options) (map[int64][]*File, error) {
	res := make(map[int64][]*File)
	for fsysIdx, fsys := range opts.fss {
		wdf := newWalkDirFunc(opts, res, fsysIdx)
		err := fs.WalkDir(fsys, ".", wdf)
		if err != nil {
			return nil, errors.Wrap(err, "walk dir")
		}
	}
	return res, nil
}

func newWalkDirFunc(opts *options, res map[int64][]*File, fsysIdx int) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if opts.errorHandler != nil {
				err = errors.Wrap(err, "walk dir")
				opts.errorHandler(err)
				return nil
			}
			return errors.Wrap(err, "")
		}
		if !d.Type().IsRegular() {
			return nil
		}
		fi, err := d.Info()
		if err != nil {
			err = errors.Wrap(err, "info")
			if opts.errorHandler != nil {
				opts.errorHandler(err)
				return nil
			}
			return err
		}
		size := fi.Size()
		if size < opts.minSize {
			return nil
		}
		fp := &File{
			FSIndex: fsysIdx,
			Path:    path,
		}
		res[size] = append(res[size], fp)
		return nil
	}
}

func getFilesByHash(opts *options, files []*File) (map[string][]*File, error) {
	res := make(map[string][]*File)
	for _, file := range files {
		h, err := hashFile(opts, file)
		if err != nil {
			if opts.errorHandler != nil {
				opts.errorHandler(err)
				continue
			}
			return nil, err
		}
		res[h] = append(res[h], file)
	}
	return res, nil
}

func hashFile(opts *options, fp *File) (string, error) {
	fsys := opts.fss[fp.FSIndex]
	f, err := fsys.Open(fp.Path)
	if err != nil {
		return "", errors.Wrap(err, "open")
	}
	defer f.Close() //nolint:errcheck
	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", errors.Wrap(err, "copy")
	}
	hs := h.Sum(nil)
	res := hex.EncodeToString(hs)
	return res, nil
}
