// +build !windows

package osfs

import (
	"syscall"

	"github.com/elastic/go-txfile/internal/vfs"
)

func sysErrKind(err error) vfs.Kind {
	err = underlyingError(err)
	switch err {
	case syscall.EDQUOT, syscall.ENOSPC, syscall.ENFILE:
		return vfs.ErrNoSpace

	case syscall.EMFILE:
		return vfs.ErrFDLimit

	case syscall.ENOTDIR:
		return vfs.ErrResolvePath

	case syscall.ENOTSUP:
		return vfs.ErrNotSupported

	case syscall.EIO:
		return vfs.ErrIO

	case syscall.EDEADLK:
		return vfs.ErrLockFailed
	}

	return vfs.ErrOSOther
}
