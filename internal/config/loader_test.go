// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseConfig(t *testing.T) {
	cfg, err := Parse(strings.NewReader(`
binary = "claude"

[profiles.personal]
env_files = ["personal.env", "common.env"]
args = ["--verbose"]

[profiles.personal.env]
ANTHROPIC_API_KEY = "test-key"
CLAUDE_CONFIG_DIR = "$HOME/.claude-personal"
`))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if cfg.Binary != "claude" {
		t.Fatalf("Binary = %q, want claude", cfg.Binary)
	}

	profile, ok := cfg.Profiles["personal"]
	if !ok {
		t.Fatal("missing personal profile")
	}

	if got, want := strings.Join(profile.EnvFiles, ","), "personal.env,common.env"; got != want {
		t.Fatalf("EnvFiles = %q, want %q", got, want)
	}
	if got, want := strings.Join(profile.Args, ","), "--verbose"; got != want {
		t.Fatalf("Args = %q, want %q", got, want)
	}
	if got, want := profile.Env["ANTHROPIC_API_KEY"], "test-key"; got != want {
		t.Fatalf("ANTHROPIC_API_KEY = %q, want %q", got, want)
	}
}

func TestWriteConfigSortsAndRoundTrips(t *testing.T) {
	cfg := New()
	cfg.Profiles["work"] = Profile{
		EnvFiles: []string{"work.env"},
		Args:     []string{"--verbose", "--label=a,b"},
		Env: map[string]string{
			"ZED": "last",
			"KEY": `quoted "value"`,
		},
	}
	cfg.Profiles["personal"] = Profile{EnvFiles: []string{"personal.env"}}

	var buf bytes.Buffer
	Write(&buf, cfg)

	got := buf.String()
	if !strings.Contains(got, "\n[profiles.personal]\n") {
		t.Fatalf("missing personal profile in output:\n%s", got)
	}
	if strings.Index(got, "[profiles.personal]") > strings.Index(got, "[profiles.work]") {
		t.Fatalf("profiles are not sorted:\n%s", got)
	}

	parsed, err := Parse(strings.NewReader(got))
	if err != nil {
		t.Fatalf("Parse(Write(cfg)) error = %v", err)
	}
	if got, want := parsed.Profiles["work"].Env["KEY"], `quoted "value"`; got != want {
		t.Fatalf("round-tripped KEY = %q, want %q", got, want)
	}
	if got, want := strings.Join(parsed.Profiles["work"].Args, "|"), "--verbose|--label=a,b"; got != want {
		t.Fatalf("round-tripped Args = %q, want %q", got, want)
	}
}

func TestSaveFileCreatesParentDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.toml")
	cfg := New()
	cfg.Profiles["personal"] = Profile{EnvFiles: []string{"personal.env"}}

	if err := SaveFile(path, cfg); err != nil {
		t.Fatalf("SaveFile() error = %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("saved config does not exist: %v", err)
	}
}
