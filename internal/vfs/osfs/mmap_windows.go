// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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
