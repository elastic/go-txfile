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
func Op(op string) *E                     { return &E{op: op} }
func Of(k error) *E                       { return &E{kind: k} }
func Msg(s string) *E                     { return &E{message: s} }
func Msgf(s string, vs ...interface{}) *E { return (&E{}).Msgf(s, vs...) }
func Wrap(err error) *E                   { return &E{cause: err} }

// chaining modifiers

func (e *E) Op(op string) *E                     { e.op = op; return e }
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
