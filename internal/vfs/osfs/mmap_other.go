// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package osfs

import "golang.org/x/sys/unix"

type mmapState struct{}

func (f *File) MMap(sz int) ([]byte, error) {
	return unix.Mmap(int(f.Fd()), 0, int(sz), unix.PROT_READ, unix.MAP_SHARED)
}

func (f *File) MUnmap(b []byte) error {
	return unix.Munmap(b)
}
