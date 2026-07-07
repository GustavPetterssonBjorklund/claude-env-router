// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package config

import (
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
