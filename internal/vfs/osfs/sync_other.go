// +build dragonfly freebsd netbsd openbsd solaris

package osfs

import (
	"golang.org/x/sys/unix"

	"github.com/elastic/go-txfile/internal/vfs"
)

// Sync uses fsync, for flushing and syncing a file to disk.  If the OS, file
// system, or disk drivers do not enforce a flush on all the intermediate
// caches and the drive itself, there is a chance of data loss and file
// corruption on power failure.
func (f *File) Sync(flags vfs.SyncFlag) error {
	// best effort
	for {
		err := f.File.Sync()
		if err == nil || (err != unix.EINTR && err != unix.EAGAIN) {
			return err
		}
	}
}
