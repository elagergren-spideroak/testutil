// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// TestInlining checks whether the methods and/or functions
// described by wantInlinable found in pkgPath are still able to
// be inlined.
func TestInlining(t *testing.T, pkgPath string, wantInlinable ...string) {
	mustHaveGoBuild(t)
	t.Parallel()
	var exe string
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}
	out, err := exec.Command(
		filepath.Join(runtime.GOROOT(), "bin", "go"+exe),
		"build",
		"--gcflags=-m",
		pkgPath,
	).CombinedOutput()
	if err != nil {
		t.Fatalf("go build: %v, %s", err, out)
	}
	got := make(map[string]bool)
	regexp.MustCompile(` can inline (\S+)`).ReplaceAllFunc(out, func(match []byte) []byte {
		got[strings.TrimPrefix(string(match), " can inline ")] = true
		return nil
	})

	for _, want := range wantInlinable {
		if !got[want] {
			t.Errorf("%q is no longer inlinable", want)
			continue
		}
		delete(got, want)
	}
	for sym := range got {
		if strings.Contains(sym, ".func") {
			continue
		}
		t.Logf("not in expected set, but also inlinable: %q", sym)
	}
}

// mustHaveGoBuild checks that the current system can build
// programs with "go build" and then run them with
// os.StartProcess or exec.Command. If not, MustHaveGoBuild calls
// t.Skip with an explanation.
func mustHaveGoBuild(t testing.TB) {
	if os.Getenv("GO_GCFLAGS") != "" {
		t.Skipf("skipping test: 'go build' not compatible with setting $GO_GCFLAGS")
	}
	if !hasGoBuild() {
		t.Skipf("skipping test: 'go build' not available on %s/%s", runtime.GOOS, runtime.GOARCH)
	}
}

// hasGoBuild reports whether the current system can build
// programs with "go build" and then run them with
// os.StartProcess or exec.Command.
func hasGoBuild() bool {
	if os.Getenv("GO_GCFLAGS") != "" {
		// It's too much work to require every caller of the go
		// command to pass along
		// "-gcflags="+os.Getenv("GO_GCFLAGS"). For now, if
		// $GO_GCFLAGS is set, report that we simply can't run go
		// build.
		return false
	}
	switch runtime.GOOS {
	case "android", "js", "ios":
		return false
	}
	return true
}
