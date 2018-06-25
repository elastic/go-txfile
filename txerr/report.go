package txerr

import (
	"bufio"
	"strings"
)

func GetMessage(in error) string {
	var msg string
	FindErrWith(in, func(in error) bool {
		if err, ok := in.(withMessage); ok {
			msg = err.Message()
			if msg != "" {
				return true
			}
		}
		return false
	})

	if msg == "" && in != nil {
		return in.Error()
	}
	return msg
}

func Report(in error) string {
	buf := &strings.Builder{}
	putStr(buf, directOp(in))
	putKind(buf, directKind(in))
	putStr(buf, directMsg(in))

	switch err := in.(type) {
	case withChild:
		putErr(buf, err.Cause())

	case withChildren:
		for _, sub := range err.Causes() {
			putSubErr(buf, sub)
		}
	}

	if buf.Len() == 0 {
		return "unknown error"
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
	if err == nil {
		return
	}

	s := err.Error()
	if s == "" {
		return
	}

	pad(b, ":\n\t")
	b.WriteString(s)
}

func putSubErr(b *strings.Builder, err error) {
	if err == nil {
		return
	}

	s := err.Error()
	if s == "" {
		return
	}

	pad(b, ":\n\t")

	// iterate lines
	r := strings.NewReader(s)
	scanner := bufio.NewScanner(r)
	first := true
	for scanner.Scan() {
		if !first {
			pad(b, "\n\t")
		} else {
			first = false
		}

		b.WriteString(scanner.Text())
	}
}

func putKind(b *strings.Builder, err error) {
	if err != nil {
		putStr(b, err.Error())
	}
}
