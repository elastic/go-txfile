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
	"bufio"
	"strings"

	"github.com/elastic/go-txfile/internal/strbld"
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
	buf := &strbld.Builder{}
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

func pad(b *strbld.Builder, p string) {
	if b.Len() > 0 {
		b.WriteString(p)
	}
}

func putStr(b *strbld.Builder, s string) {
	if s != "" {
		pad(b, ": ")
		b.WriteString(s)
	}
}

func putErr(b *strbld.Builder, err error) {
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

func putSubErr(b *strbld.Builder, err error) {
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

func putKind(b *strbld.Builder, err error) {
	if err != nil {
		putStr(b, err.Error())
	}
}
