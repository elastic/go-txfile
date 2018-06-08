package osfs

import (
	"os"
	"syscall"
)

// File implements vfs.File for the current target operating system.
type File struct {
	*os.File
	state osFileState
}

type osFileState struct {
	mmap mmapState
	lock lockState
	sync syncState
}

var errno0 = syscall.Errno(0)

func Open(path string, mode os.FileMode) (*File, error) {
	flags := os.O_RDWR | os.O_CREATE
	f, err := os.OpenFile(path, flags, mode)
	return &File{File: f}, err
}

func (o *File) Size() (int64, error) {
	stat, err := o.File.Stat()
	if err != nil {
		return -1, err
	}
	return stat.Size(), nil
}

func normalizeSysError(err error) error {
	if err == nil || err == errno0 {
		return nil
	}
	return err
}
