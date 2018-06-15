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
	Truncate(int64) error

	Lock(exclusive, blocking bool) error
	Unlock() error

	MMap(sz int) ([]byte, error)
	MUnmap([]byte) error

	// If a write/flush fails due to IO errors or the disk running out of space,
	// the kernel internally marks the error on the 'page'. Fsync will finally
	// return the error, but reset the error on failed writes. Subsequent fsync
	// operations will not report errors for former failed pages, even if the
	// pages are not written again. Therefore, if fsync fails, we must assume all
	// write operations - since the last successfull fsync - have failed and
	// reinitiate all writes.
	// According to [1] Linux, OpenBSD, and NetBSD are known to silently clear
	// errors on fsync fail.
	//
	// [1]: https://lwn.net/Articles/752098/
	Sync(flags SyncFlag) error
}

type SyncFlag uint8

const (
	SyncAll SyncFlag = 0

	// SyncDataOnly will only flush the file data, without enforcing an update on
	// the file metadata (like file size or modification time).
	SyncDataOnly SyncFlag = 1 << iota
)
