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
	"github.com/elastic/go-txfile/internal/vfs"
	"github.com/elastic/go-txfile/txerr"
)

type reason interface {
	txerr.ErrorBuild
}

type ErrKind int

// internal txfile error kinds

//go:generate stringer -type=ErrKind -linecomment=true
const (
	InternalError      ErrKind = iota // internal error
	FileCreationFailed                // can not create file
	InvalidConfig                     // configuration error
	InvalidFileSize                   // invalid file size
	InvalidMetaPage                   // meta page invalid
	InvalidOp                         // invalid operation
	InvalidPageID                     // page id out of bounds
	InvalidParam                      // invalid parameter
	OutOfMemory                       // out of memory
	TxCommitFail                      // transaction failed during commit
	TxFailed                          // transaction failed
	TxFinished                        // finished transaction
	TxReadOnly                        // readonly transaction
	endOfErrKind                      // unknown error kind
)

// re-export file system error kinds (from internal/vfs)

const (
	PermissionError       = vfs.ErrPermission
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

func (k ErrKind) Error() string {
	if k > endOfErrKind {
		k = endOfErrKind
	}
	return k.String()
}

func raiseOutOfBounds(op string, id PageID) *txerr.E {
	return txerr.Op(op).Of(InvalidPageID).Msgf("out of bounds page id %v", id)
}
