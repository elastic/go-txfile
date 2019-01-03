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
	"testing"

	txfile "github.com/elastic/go-txfile"
	"github.com/elastic/go-txfile/internal/mint"
)

type testObserveLast statEntry

type statEntry struct {
	kind statKind
	off  uintptr

	// OnInit statis
	version   uint32
	available uint

	// per callback stats
	flush FlushStats
	read  ReadStats
	ack   ACKStats
}

type statKind uint8

const (
	statNone statKind = iota
	statOnOpen
	statOnFlush
	statOnRead
	statOnACK
)

var isWindows bool

func init() {
	isWindows = runtime.OS == "windows"
}

func TestObserveStats(testing *testing.T) {
	t := mint.NewWith(testing, func(sub *mint.T) func() {
		pushTracer(mint.NewTestLogTracer(sub, logTracer))
		return popTracer
	})

	withQueue := func(fn func(*mint.T, *testQueue, *statEntry)) func(*mint.T) {
		return func(t *mint.T) {
			stat := &statEntry{}
			qu, teardown := setupQueue(t, config{
				File: txfile.Options{
					MaxSize:  128 * 1024, // default file size of 128 pages
					PageSize: 1024,
				},
				Queue: Settings{
					WriteBuffer: defaultMinPages, // buffer up to 5 pages
					Observer:    newTestObserver(stat),
				},
			})
			defer teardown()

			stat.reset()
			fn(t, qu, stat)
		}
	}

	t.Run("open empty queue", withQueue(func(t *mint.T, qu *testQueue, stat *statEntry) {
		qu.Reopen()
		t.Equal(stat.kind, statOnOpen)
		t.Equal(int(stat.version), queueVersion) // current version
		t.Equal(int(stat.available), 0)
	}))

	t.Run("open non-empty queue", withQueue(func(t *mint.T, qu *testQueue, stat *statEntry) {
		// write 3 events
		qu.append("a", "b", "c")
		qu.flush()

		// validate
		qu.Reopen()
		t.Equal(stat.kind, statOnOpen)
		t.Equal(int(stat.version), queueVersion) // current version
		t.Equal(int(stat.available), 3)
	}))

	t.Run("small write with explicit flush", withQueue(func(t *mint.T, qu *testQueue, stat *statEntry) {
		qu.append("a", "bc", "def")
		t.Equal(statNone, stat.kind, "unexpected stats update")

		qu.flush()
		t.Equal(statOnFlush, stat.kind)
		t.Equal(FlushStats{
			// do not compare duration and timestamps ;)
			Duration: stat.flush.Duration,
			Oldest:   stat.flush.Oldest,
			Newest:   stat.flush.Newest,

			// validated fields
			Failed:      false,
			OutOfMemory: false,
			Pages:       1,
			Allocate:    1,
			Events:      3,
			BytesTotal:  6,
			BytesMin:    1,
			BytesMax:    3,
		}, stat.flush)
		t.True(stat.flush.Duration > 0, "flush duration should be > 0")
		t.False(stat.flush.Oldest.IsZero(), "oldest timestamp must not be 0")
		t.False(stat.flush.Newest.IsZero(), "newest timestamp must not be 0")

		if isWindows {
			t.True(stat.flush.Oldest != stat.flush.Newest, "timestamps do match")
		}
	}))

	t.Run("big write with implcicit flush", withQueue(func(t *mint.T, qu *testQueue, stat *statEntry) {
		var msg [5000]byte
		qu.append(string(msg[:]))

		t.Equal(statOnFlush, stat.kind)
		t.Equal(FlushStats{
			// do not compare duration and timestamps ;)
			Duration: stat.flush.Duration,
			Oldest:   stat.flush.Oldest,
			Newest:   stat.flush.Newest,

			// validated fields
			Failed:      false,
			OutOfMemory: false,
			Pages:       6,
			Allocate:    6,
			Events:      1,
			BytesTotal:  5000,
			BytesMin:    5000,
			BytesMax:    5000,
		}, stat.flush)
		t.True(stat.flush.Duration > 0, "flush duration should be > 0")
		t.False(stat.flush.Oldest.IsZero(), "oldest timestamp must not be 0")
		t.False(stat.flush.Newest.IsZero(), "newest timestamp must not be 0")
		t.True(stat.flush.Oldest == stat.flush.Newest, "timestamps do not match")
	}))

	t.Run("flush on close", withQueue(func(t *mint.T, qu *testQueue, stat *statEntry) {
		qu.append("a", "bc", "def")
		qu.Close()

		t.Equal(statOnFlush, stat.kind)
		t.Equal(FlushStats{
			// do not compare duration and timestamps ;)
			Duration: stat.flush.Duration,
			Oldest:   stat.flush.Oldest,
			Newest:   stat.flush.Newest,

			// validated fields
			Failed:      false,
			OutOfMemory: false,
			Pages:       1,
			Allocate:    1,
			Events:      3,
			BytesTotal:  6,
			BytesMin:    1,
			BytesMax:    3,
		}, stat.flush)
		t.True(stat.flush.Duration > 0, "flush duration should be > 0")
		t.False(stat.flush.Oldest.IsZero(), "oldest timestamp must not be 0")
		t.False(stat.flush.Newest.IsZero(), "newest timestamp must not be 0")

		if !isWindows {
			t.True(stat.flush.Oldest != stat.flush.Newest, "timestamps do match")
		}
	}))
}

func newTestObserver(stat *statEntry) *testObserveLast {
	return (*testObserveLast)(stat)
}

func (t *testObserveLast) OnQueueInit(off uintptr, version uint32, available uint) {
	t.set(off, statOnOpen, statEntry{version: version, available: available})
}

func (t *testObserveLast) OnQueueFlush(off uintptr, stats FlushStats) {
	t.set(off, statOnFlush, statEntry{flush: stats})
}

func (t *testObserveLast) OnQueueRead(off uintptr, stats ReadStats) {
	t.set(off, statOnRead, statEntry{read: stats})
}

func (t *testObserveLast) OnQueueACK(off uintptr, stats ACKStats) {
	t.set(off, statOnACK, statEntry{ack: stats})
}

func (t *testObserveLast) set(off uintptr, kind statKind, e statEntry) {
	*t = testObserveLast(e)
	t.off, t.kind = off, kind
}

func (s *statEntry) reset() { *s = statEntry{} }
