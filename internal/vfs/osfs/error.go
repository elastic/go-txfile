package osfs

import (
	"os"

	"github.com/elastic/go-txfile/internal/vfs"
)

func errKind(err error) vfs.Kind {
	if os.IsPermission(err) {
		return vfs.ErrPermissions
	}
	if os.IsExist(err) {
		return vfs.ErrExist
	}
	if os.IsNotExist(err) {
		return vfs.ErrExist
	}

	switch err {
	case os.ErrClosed:
		return vfs.ErrClosed
	default:
		return sysErrKind(err)
	}
}

func normalizeSysError(err error) error {
	err = underlyingError(err)
	if err == nil || err == errno0 {
		return nil
	}
	return err
}

func underlyingError(in error) error {
	switch err := in.(type) {
	case *os.PathError:
		return err.Err

	case *os.LinkError:
		return err.Err

	case *os.SyscallError:
		return err.Err

	default:
		return err
	}
}
