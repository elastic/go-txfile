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

package osfs

import (
	"fmt"

	flock "github.com/theckman/go-flock"
)

const (
	lockExt = ".lock"
)

type lockState struct {
	*flock.Flock
}

func (f *File) Lock(exclusive, blocking bool) error {
	if f.state.lock.Flock != nil {
		return fmt.Errorf("file %v is already locked", f.Name())
	}

	var ok bool
	var err error
	lock := flock.NewFlock(f.Name() + lockExt)
	if blocking {
		err = lock.Lock()
		ok = err == nil
	} else {
		ok, err = lock.TryLock()
	}

	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("file %v can not be locked right now", f.Name())
	}

	f.state.lock.Flock = lock
	return nil
}

func (f *File) Unlock() error {
	if f.state.lock.Flock == nil {
		return fmt.Errorf("file %v is not locked", f.Name())
	}

	err := f.state.lock.Unlock()
	if err == nil {
		f.state.lock.Flock = nil
	}
	return err
}
