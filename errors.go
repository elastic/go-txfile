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

package txfile

import (
	"errors"

	"github.com/elastic/go-txfile/internal/vfs"
	"github.com/elastic/go-txfile/txerr"
)

type reason interface {
	txerr.ErrorBuild
}

type ErrKind int

// txfile internal error kinds
var (
// InvalidConfig = txerr.NewKind("configuration error")
)

// file system error kinds (map internal/vfs errors)
const (
	PermissionError       = vfs.ErrPermissions
	FileExists            = vfs.ErrExist
	FileDoesNotExist      = vfs.ErrNotExist
	FileClosed            = vfs.ErrClosed
	NoDiskSpace           = vfs.ErrNoSpace
	FDLimit               = vfs.ErrFDLimit
	CantResolvePath       = vfs.ErrResolvePath
	IOError               = vfs.ErrIO
	OSOtherError          = vfs.ErrOSOther
	OperationNotSupported = vfs.ErrNotSupported
	LockFailed            = vfs.ErrLockFailed
)

// internal txfile error kinds
const (
	InvalidConfig ErrKind = iota
	InvalidParam
	InvalidFileSize
	InvalidPageID
	FileCreationFailed
	TxFailed
)

var kindStr = [...]string{
	"configuration error",
	"invalid parameter",
	"invalid file size",
	"page id out of bounds",
	"can not create file",
	"transaction failed",
}

func (k ErrKind) Error() string {
	if k > 0 && int(k) < len(kindStr) {
		return kindStr[k]
	}
	return "unknown error kind"
}

// TODO: obsolete errors:
var (
	// settings errors
	errReadOnlyUpdateSize = errors.New("can not update the file size in read only mode")

	// file meta page validation errors

	errMagic    = errors.New("invalid magic number")
	errVersion  = errors.New("invalid version number")
	errChecksum = errors.New("checksum mismatch")

	// file sizing errors

	// errMmapTooLarge    = errors.New("mmap too large")
	errFileSizeTooLage = errors.New("max file size to large for this system")
	// errInvalidFileSize = errors.New("invalid file size")

	// page access/allocation errors

	errOutOfBounds   = errors.New("out of bounds page id")
	errOutOfMemory   = errors.New("out of memory")
	errFreedPage     = errors.New("trying to access an already freed page")
	errPageFlushed   = errors.New("page is already flushed")
	errTooManyBytes  = errors.New("contents exceeds page size")
	errNoPageData    = errors.New("accessing page without contents")
	errFreeDirtyPage = errors.New("freeing dirty page")

	// transaction errors

	errTxFinished = errors.New("transaction has already been closed")
	errTxReadonly = errors.New("readonly transaction")
)
