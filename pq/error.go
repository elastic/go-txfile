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

package pq

import (
	"fmt"

	"github.com/elastic/go-txfile"

	"github.com/elastic/go-txfile/internal/strbld"
	"github.com/elastic/go-txfile/txerr"
)

// ErrKind provides the pq related error kinds
type ErrKind int

type reason interface {
	txerr.Error
}

type Error struct {
	op    string
	kind  error
	cause error
	ctx   errorCtx
	msg   string
}

type errorCtx struct {
	id queueID

	isPage bool
	page   txfile.PageID
}

//go:generate stringer -type=ErrKind -linecomment=true
const (
	NoError          ErrKind = iota // no error
	InitFailed                      // failed to initialize queue
	InvalidParam                    // invalid parameter
	InvalidPageSize                 // invalid page size
	InvalidConfig                   // invalid queue config
	QueueClosed                     // queue is already closed
	ReaderClosed                    // reader is already closed
	WriterClosed                    // writer is already closed
	NoQueueRoot                     // no queue root
	InvalidQueueRoot                // queue root is invalid
	QueueVersion                    // unsupported queue version
	ACKEmptyQueue                   // invalid ack on empty queue
	ACKTooMany                      // too many events acked
	SeekFail                        // failed to seek to next page
	ReadFail                        // failed to read page
)

func (k ErrKind) Error() string {
	return k.String()
}

func (e *Error) Error() string              { return txerr.Report(e, false) }
func (e *Error) Format(s fmt.State, c rune) { txerr.Format(e, s, c) }
func (e *Error) Op() string                 { return e.op }
func (e *Error) Kind() error                { return e.kind }
func (e *Error) Cause() error               { return e.cause }
func (e *Error) Context() string            { return e.ctx.String() }
func (e *Error) Message() string            { return e.msg }

func (ctx *errorCtx) String() string {
	buf := &strbld.Builder{}
	if ctx.id != 0 {
		buf.Fmt("queueID=%v", ctx.id)
	}

	if ctx.isPage {
		buf.Pad(" ")
		buf.Fmt("page=%v", ctx.page)
	}
	return buf.String()
}

func IsQueueCorrupt(err error) bool {
	for _, kind := range []ErrKind{InvalidQueueRoot, SeekFail} {
		if txerr.Is(kind, err) {
			return true
		}
	}
	return false
}

func errOp(op string) *Error {
	return &Error{op: op}
}

func wrapErr(op string, cause error) *Error {
	return errOp(op).causedBy(cause)
}

func (e *Error) of(k ErrKind) *Error {
	e.kind = k
	return e
}

func (e *Error) causedBy(cause error) *Error {
	e.cause = cause
	other, ok := cause.(*Error)
	if !ok {
		return e
	}

	// merge error and cause context such that the cause context only reports
	// fields that differ from the current context.

	errCtx := &e.ctx
	causeCtx := &other.ctx

	if errCtx.id == causeCtx.id {
		causeCtx.id = 0 // delete common queue id from cause context
	}
	if errCtx.isPage && causeCtx.isPage && errCtx.page == causeCtx.page {
		causeCtx.isPage = false // delete common page id from cause context
	}

	return e
}

func (e *Error) report(m string) *Error {
	e.msg = m
	return e
}
