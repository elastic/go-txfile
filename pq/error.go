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

import "errors"

var (
	errNODelegate          = errors.New("delegate must not be nil")
	errInvalidPagesize     = errors.New("invalid page size")
	errClosed              = errors.New("queue closed")
	errNoQueueRoot         = errors.New("no queue root")
	errIncompleteQueueRoot = errors.New("incomplete queue root")
	errInvalidVersion      = errors.New("invalid queue version")
	errACKEmptyQueue       = errors.New("ack on empty queue")
	errACKTooManyEvents    = errors.New("too many events have been acked")
	errSeekPageFailed      = errors.New("failed to seek to next page")
)
