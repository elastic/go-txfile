package osfs

import (
	"golang.org/x/sys/unix"

	"github.com/elastic/go-txfile/internal/vfs"
)

type syncState struct {
	noFullSync bool
}

// Sync uses fnctl or fsync in order to flush the file buffers to disk.
// According to the darwin fsync man page[1], usage of sync is not safe. On
// darwin, fsync will only flush the OS file cache to disk, but this won't
// enforce a cache flush on the drive itself. Without forcing the cache flush,
// writes can still be out of order or get lost on power failure.
// According to the man page[1] fcntl with F_FULLFSYNC[2] is required. F_FULLFSYNC
// might not be supported for the current file system. In this case we will
// fallback to fsync.
//
// [1]: https://www.unix.com/man-page/osx/2/fsync
// [2]: https://www.unix.com/man-page/osx/2/fcntl
func (f *File) Sync(flags vfs.SyncFlag) error {
	if f.state.sync.noFullSync {
		return f.syncWithFSync()
	}
	return f.syncWithFcntl()
}

func (f *File) syncWithFcntl() error {
	for {
		_, err := unix.FcntlInt(f.File.Fd(), unix.F_FULLFSYNC, 0)
		err = normalizeSysError(err)
		if err == nil {
			return nil
		}

		switch err {
		// try again
		case unix.EINTR:
		case unix.EAGAIN:

		// fallback to fsync?
		// XXX: always fallback to fsync in future calls?
		case unix.EINVAL:
			f.state.sync.noFullSync = true
			return f.syncWithFSync()

		default:
			return err
		}
	}
}

func (f *File) syncWithFSync() error {
	for {
		err := f.File.Sync()
		if err != unix.EINTR && err != unix.EAGAIN {
			return err
		}
	}
}
