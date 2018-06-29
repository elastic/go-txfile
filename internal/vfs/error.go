package vfs

import "strings"

type Error struct {
	op   string
	kind error
	path string
	err  error
}

type Kind int

const (
	ErrOSOther Kind = iota
	ErrPermissions
	ErrExist
	ErrNotExist
	ErrClosed
	ErrNoSpace
	ErrFDLimit
	ErrResolvePath
	ErrIO
	ErrNotSupported
	ErrLockFailed
	ErrUnlockFailed
)

var kindStr = [...]string{
	"",
	"permission denied",
	"file already exists",
	"file does not exist",
	"file already closed",
	"no space or quota exhausted",
	"process file desciptor limit reached",
	"cannot resolve path",
	"read/write IO error",
	"operation not supported",
	"file lock failed",
	"file unlock failed",
}

func (k Kind) Error() string {
	if k > 0 && int(k) < len(kindStr) {
		return kindStr[k]
	}
	return "unknown"
}

func Err(op string, kind Kind, path string, err error) *Error {
	return &Error{op: op, kind: kind, path: path, err: err}
}

func (e *Error) Op() string   { return e.op }
func (e *Error) Kind() error  { return e.kind }
func (e *Error) Path() string { return e.path }
func (e *Error) Cause() error { return e.err }

func (e *Error) Error() string {
	buf := &strings.Builder{}
	putStr(buf, e.op)
	putErr(buf, e.kind)
	putStr(buf, e.path)
	putErr(buf, e.err)

	if buf.Len() == 0 {
		return "no error"
	}
	return buf.String()
}

func pad(b *strings.Builder, p string) {
	if b.Len() > 0 {
		b.WriteString(p)
	}
}

func putStr(b *strings.Builder, s string) {
	if s != "" {
		pad(b, ": ")
		b.WriteString(s)
	}
}

func putErr(b *strings.Builder, err error) {
	if err != nil {
		putStr(b, err.Error())
	}
}
