package osfs

import (
	"golang.org/x/sys/unix"

	"github.com/elastic/go-txfile/internal/vfs"
)

type syncState struct {
	noDataOnly bool
}

// Sync uses fsync or fdatasync (if vfs.SyncDataOnly flag is set).
//
// Handling write-back errors is at a mess in older linux kernels [1].
// With mixed read-write operations, there is a chance that write-back errors
// are never reported to user-space applications, as error flags are cleared in
// the caches.
// Error handling was somewhat improved in 4.13 [2][3], such that errors will
// actually be reported on fsync (more improvements have been added to 4.16).
//
// [1]: https://lwn.net/Articles/718734//
// [2]: https://lwn.net/Articles/724307/
// [3]: https://lwn.net/Articles/724232/
func (f *File) Sync(flags vfs.SyncFlag) error {
	dataOnly := (flags & vfs.SyncDataOnly) != 0
	for {
		err := f.doSync(!f.state.sync.noDataOnly && dataOnly)
		if err == nil || (err != unix.EINTR && err != unix.EAGAIN) {
			return err
		}
	}
}

func (f *File) doSync(dataOnly bool) error {
	if dataOnly {
		err := normalizeSysError(unix.Fdatasync(int(f.File.Fd())))
		if err == unix.ENOSYS {
			f.state.sync.noDataOnly = true
			return f.File.Sync()
		}
	}
	return f.File.Sync()
}
