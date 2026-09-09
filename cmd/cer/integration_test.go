// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestCERExecutesConfiguredProfile exercises the installed-user path through
// the actual cer executable: argument parsing, config loading, dotenv loading,
// environment merging, and child-process execution.
func TestCERExecutesConfiguredProfile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("integration fixture uses a POSIX shell")
	}

	dir := t.TempDir()
	cerPath := filepath.Join(dir, "cer")
	build := exec.Command("go", "build", "-o", cerPath, ".")
	build.Dir = "."
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, output)
	}

	fakeClaude := filepath.Join(dir, "fake-claude")
	writeExecutable(t, fakeClaude, `#!/bin/sh
set -eu
printf 'args:'
for arg do printf ' <%s>' "$arg"; done
printf '\n'
printf 'CER_TEST_FILE=%s\n' "${CER_TEST_FILE-}"
printf 'CER_TEST_INLINE=%s\n' "${CER_TEST_INLINE-}"
printf 'CER_TEST_PROCESS=%s\n' "${CER_TEST_PROCESS-}"
`)

	configPath := filepath.Join(dir, "config.toml")
	writeFile(t, filepath.Join(dir, "profile.env"), "CER_TEST_FILE=from-file\n")
	writeFile(t, configPath, `binary = "`+fakeClaude+`"

[profiles.integration]
env_files = ["profile.env"]
args = ["--configured", "value"]

[profiles.integration.env]
CER_TEST_INLINE = "from-inline"
`)
	t.Setenv("CER_TEST_PROCESS", "from-process")

	cmd := exec.Command(cerPath, "--config", configPath, "integration", "--", "--prompt", "hello world")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cer: %v\n%s", err, output)
	}

	got := string(output)
	for _, want := range []string{
		"args: <--configured> <value> <--prompt> <hello world>",
		"CER_TEST_FILE=from-file",
		"CER_TEST_INLINE=from-inline",
		"CER_TEST_PROCESS=from-process",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("cer output = %q, want it to contain %q", got, want)
		}
	}
}

func writeFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func writeExecutable(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o700); err != nil {
		t.Fatalf("write executable %s: %v", path, err)
	}
}
