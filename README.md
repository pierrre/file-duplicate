# File Duplicate

[![Go Reference](https://pkg.go.dev/badge/github.com/pierrre/file-duplicate.svg)](https://pkg.go.dev/github.com/pierrre/file-duplicate)

## Features

- Find duplicate files
- Command line ready to use
- Library that can be integrated in a project

## Usage

```bash
# Local build
make build
./build/file-duplicate -h

# Remote install
go install github.com/pierrre/file-duplicate/cmd/file-duplicate@latest

# Module install
go get github.com/pierrre/file-duplicate@latest
```

## Implementation

- Walk the filesystems
- Group files by identical size
- Compute the SHA256 hash of same size files
- Files with same hash are duplicates
