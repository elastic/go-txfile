package osfs

import "os"

// File implements vfs.File for the current target operating system.
type File struct {
	*os.File
	state osFileState
}

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
