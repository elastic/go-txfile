package vfs

import (
	"io"
)

type File interface {
	io.Closer
	io.WriterAt
	io.ReaderAt

	Name() string
	Size() (int64, error)
	Sync() error
	Truncate(int64) error

	Lock(exclusive, blocking bool) error
	Unlock() error

	MMap(sz int) ([]byte, error)
	MUnmap([]byte) error
}
