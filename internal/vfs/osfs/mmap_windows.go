package osfs

import (
	"os"
	"reflect"
	"unsafe"

	"golang.org/x/sys/windows"
)

type mmapState struct {
	windows.Handle
}

func (f *File) MMap(sz int) ([]byte, error) {
	szHi, szLo := uint32(sz>>32), uint32(sz)
	hdl, err := windows.CreateFileMapping(windows.Handle(f.Fd()), nil, windows.PAGE_READONLY, szHi, szLo, nil)
	if hdl == 0 {
		return nil, os.NewSyscallError("CreateFileMapping", err)
	}

	// map memory
	addr, err := windows.MapViewOfFile(hdl, windows.FILE_MAP_READ, 0, 0, uintptr(sz))
	if addr == 0 {
		windows.CloseHandle(hdl)
		return nil, os.NewSyscallError("MapViewOfFile", err)
	}

	f.state.mmap.Handle = hdl

	slice := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(addr),
		Len:  sz,
		Cap:  sz}))
	return slice, nil
}

func (f *File) MUnmap(b []byte) error {
	err1 := windows.UnmapViewOfFile(uintptr(unsafe.Pointer(&b[0])))
	b = nil

	err2 := windows.CloseHandle(f.state.mmap.Handle)
	f.state.mmap.Handle = 0

	if err1 != nil {
		return os.NewSyscallError("UnmapViewOfFile", err1)
	} else if err2 != nil {
		return os.NewSyscallError("CloseHandle", err2)
	}
	return nil
}
