// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package osfs

import "golang.org/x/sys/unix"

type lockState struct{}

func (f *File) Lock(exclusive, blocking bool) error {
	flags := unix.LOCK_SH
	if exclusive {
		flags = unix.LOCK_EX
	}
	if !blocking {
		flags |= unix.LOCK_NB
	}

	return unix.Flock(int(f.Fd()), flags)
}

func (f *File) Unlock() error {
	return unix.Flock(int(f.Fd()), unix.LOCK_UN)
}
