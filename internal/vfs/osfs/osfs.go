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
	"syscall"
)

// File implements vfs.File for the current target operating system.
type File struct {
	*os.File
	state osFileState
}

type osFileState struct {
	mmap mmapState
	lock lockState
	sync syncState
}

var errno0 = syscall.Errno(0)

func Open(path string, mode os.FileMode) (*File, error) {
	flags := os.O_RDWR | os.O_CREATE
	f, err := os.OpenFile(path, flags, mode)
	return &File{File: f}, err
}

func (o *File) Size() (int64, error) {
	stat, err := o.File.Stat()
	if err != nil {
		return -1, err
	}
	return stat.Size(), nil
}

func normalizeSysError(err error) error {
	if err == nil || err == errno0 {
		return nil
	}
	return err
}
