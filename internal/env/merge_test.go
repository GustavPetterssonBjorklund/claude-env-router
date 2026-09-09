// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package env

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/config"
)

func TestApplyProfilePrecedenceAndExpansion(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "base.env"), "TOKEN=from-file\nROOT=$HOME/project\n")
	writeFile(t, filepath.Join(dir, "override.env"), "TOKEN=from-second-file\n")

	values := FromEnviron([]string{"HOME=/home/test", "TOKEN=from-process"})
	err := ApplyProfile(values, dir, config.Profile{
		EnvFiles: []string{"base.env", "override.env"},
		Env: map[string]string{
			"TOKEN": "from-inline",
			"PATH":  "$ROOT/bin",
		},
	}, nil)
	if err != nil {
		t.Fatalf("ApplyProfile() error = %v", err)
	}

	want := map[string]string{
		"HOME":  "/home/test",
		"TOKEN": "from-inline",
		"ROOT":  "/home/test/project",
		"PATH":  "/home/test/project/bin",
	}
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("values = %#v, want %#v", values, want)
	}
}

func TestToEnvironSortsKeys(t *testing.T) {
	got := ToEnviron(map[string]string{"B": "2", "A": "1"})
	want := []string{"A=1", "B=2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ToEnviron() = %#v, want %#v", got, want)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
