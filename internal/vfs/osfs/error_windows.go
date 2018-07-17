package osfs

import (
	"syscall"

	"github.com/elastic/go-txfile/internal/vfs"
)

const (
	ERROR_DISK_FULL             syscall.Errno = 62
	ERROR_DISK_QUOTA_EXCEEDED   syscall.Errno = 1295
	ERROR_TOO_MANY_OPEN_FILES   syscall.Errno = 4
	ERROR_LOCK_FAILED           syscall.Errno = 167
	ERROR_CANT_RESOLVE_FILENAME syscall.Errno = 1921
)

func sysErrKind(err error) vfs.Kind {
	switch underlyingError(err) {

	case ERROR_DISK_FULL, ERROR_DISK_QUOTA_EXCEEDED:
		return vfs.ErrNoSpace

	case ERROR_TOO_MANY_OPEN_FILES:
		return vfs.ErrFDLimit

	case ERROR_LOCK_FAILED:
		return vfs.ErrLockFailed

	case ERROR_CANT_RESOLVE_FILENAME:
		return vfs.ErrResolvePath

	default:
		return vfs.ErrOSOther
	}
}
