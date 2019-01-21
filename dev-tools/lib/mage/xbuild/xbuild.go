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

package xbuild

import (
	"fmt"

	"github.com/magefile/mage/mg"
)

type Regsitry struct {
	table map[OSArch]Provider
}

type Provider interface {
	Build() error
	Run(env map[string]string, cmdAndArgs ...string) error
}

type OSArch struct {
	OS   string
	Arch string
}

func NewRegistry(tbl map[OSArch]Provider) *Regsitry {
	return &Regsitry{tbl}
}

func (r *Regsitry) Find(os, arch string) (Provider, error) {
	p := r.table[OSArch{os, arch}]
	if p == nil {
		return nil, fmt.Errorf("No provider for %v:%v defined", os, arch)
	}
	return p, nil
}

func (r *Regsitry) With(os, arch string, fn func(Provider) error) error {
	p, err := r.Find(os, arch)
	if err != nil {
		return err
	}

	mg.Deps(p.Build)
	return fn(p)
}
