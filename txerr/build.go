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

package txerr

import (
	"fmt"
	"unsafe"
)

type E struct {
	op      string
	kind    error
	message string
	cause   error
}

type errE struct{ E }

// constructors
func Op(op string) *E                     { return newE().Op(op) }
func Opf(s string, vs ...interface{}) *E  { return newE().Opf(s, vs...) }
func Of(k error) *E                       { return newE().Of(k) }
func Msg(s string) *E                     { return newE().Msg(s) }
func Msgf(s string, vs ...interface{}) *E { return newE().Msgf(s, vs...) }
func Wrap(err error) *E                   { return newE().CausedBy(err) }

// chaining modifiers

func newE() *E                                   { return &E{} }
func (e *E) Op(op string) *E                     { e.op = op; return e }
func (e *E) Opf(s string, vs ...interface{}) *E  { e.op = fmt.Sprintf(s, vs...); return e }
func (e *E) Of(kind error) *E                    { e.kind = kind; return e }
func (e *E) CausedBy(err error) *E               { e.cause = err; return e }
func (e *E) Msg(s string) *E                     { e.message = s; return e }
func (e *E) Msgf(s string, vs ...interface{}) *E { e.message = fmt.Sprintf(s, vs...); return e }

// ErrorBuild + error interface implementation

func (e *E) Error() string { return e.Err().Error() }
func (e *E) Err() Error    { return (*errE)(unsafe.Pointer(e)) }

func (e *errE) Error() string   { return Report(e) }
func (e *errE) Op() string      { return e.E.op }
func (e *errE) Kind() error     { return e.E.kind }
func (e *errE) Message() string { return e.E.message }
func (e *errE) Cause() error    { return e.E.cause }
