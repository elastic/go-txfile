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

//+build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"

	"github.com/joeshaw/multierror"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/elastic/go-txfile/dev-tools/lib/mage/gotool"
)

type Check mg.Namespace
type Prepare mg.Namespace
type Build mg.Namespace
type Info mg.Namespace

const buildHome = "build"

type envVar struct{ name, other, doc string }

// environment variables
var (
	envBuildOS    = defEnv("BUILD_OS", "", "(string) set compiler target GOOS")
	envBuildArch  = defEnv("BUILD_ARCH", "", "(string) set compiler target GOARCH")
	envTestUseBin = defEnv("TEST_USE_BIN", "", "(bool) reuse prebuild test binary when running tests")
)

var envVars = map[string]*envVar{}

func defEnv(name, value, doc string) *envVar {
	e := &envVar{name: name, other: value, doc: doc}
	if _, exists := envVars[name]; exists {
		panic(fmt.Errorf("env variable '%v' already registered", name))
	}
	envVars[name] = e
	return e
}

func (e *envVar) get() string {
	v := os.Getenv(e.name)
	if v == "" {
		return e.other
	}
	return v
}

// targets

func (Info) Env() {
	keys := make([]string, 0, len(envVars))
	for k := range envVars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := envVars[k]
		fmt.Printf("%v: %v\n", k, v.doc)
	}
}

func (Check) Lint() error {
	return errors.New("TODO: implement me")
}

func (Prepare) Dirs() error {
	return mkdir("build")
}

func (Build) Mage() error {
	mg.Deps(Prepare.Dirs)

	goos := envDefault(envBuildOS, runtime.GOOS)
	goarch := envDefault(envBuildArch, runtime.GOARCH)
	out := filepath.Join(buildHome, fmt.Sprintf("mage-%v-%v", goos, goarch))
	return sh.Run("mage", "-f", "-goos="+goos, "-goarch="+goarch, "-compile", out)
}

func (Build) Test() error {
	mg.Deps(Prepare.Dirs)

	return withList(gotool.ListProjectPackages, each, func(pkg string) error {
		tst := gotool.Test
		return tst(
			tst.OS(env(envBuildOS)),
			tst.ARCH(env(envBuildArch)),
			tst.Create(),
			tst.WithCoverage(""),
			tst.Out(path.Join(buildHome, pkg, path.Base(pkg))),
			tst.Package(pkg),
		)
	})
}

func Test() error {
	mg.Deps(Prepare.Dirs)

	return withList(gotool.ListProjectPackages, each, func(pkg string) error {
		tst := gotool.Test
		fmt.Println("Test:", pkg)
		bin := path.Join(buildHome, pkg, path.Base(pkg))
		return tst(
			tst.Use(useIf(bin, existsFile(bin) && envBool(envTestUseBin))),
			tst.WithCoverage(path.Join(buildHome, pkg, "cover.out")),
			tst.Out(bin),
			tst.Package(pkg),
			tst.Verbose(),
		)
	})
}

func Clean() error {
	return sh.Rm(buildHome)
}

// helpers

func withList(
	gen func() ([]string, error),
	mode func(...func() error) error,
	fn func(string) error,
) error {
	list, err := gen()
	if err != nil {
		return err
	}

	ops := make([]func() error, len(list))
	for i, v := range list {
		v := v
		ops[i] = func() error { return fn(v) }
	}

	return mode(ops...)
}

func useIf(s string, b bool) string {
	if b {
		return s
	}
	return ""
}

func existsFile(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode().IsRegular()
}

func env(ev *envVar) string { return ev.get() }

func envDefault(ev *envVar, other string) string {
	v := env(ev)
	if v == "" {
		return other
	}
	return v
}

func envBool(ev *envVar) bool {
	b, err := strconv.ParseBool(env(ev))
	return err == nil && b
}

func mkdirs(paths ...string) error {
	for _, p := range paths {
		if err := mkdir(p); err != nil {
			return err
		}
	}
	return nil
}

func mkdir(path string) error {
	return os.MkdirAll(path, os.ModeDir|0700)
}

func each(ops ...func() error) error {
	var errs multierror.Errors
	for _, op := range ops {
		if err := op(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs.Err()
}

func and(ops ...func() error) error {
	for _, op := range ops {
		if err := op(); err != nil {
			return err
		}
	}
	return nil
}
