// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/config"
	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/tui"
)

func TestRunNoProfileNonInteractiveErrorsClearly(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunWithInput([]string{"cer"}, strings.NewReader(""), &stdout, &stderr)

	if code != 2 {
		t.Fatalf("exit code = %d, want 2", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if got := stderr.String(); !strings.Contains(got, "no profile specified") {
		t.Fatalf("stderr = %q, want no profile error", got)
	}
}

func TestCreateProfileWritesConfigAndEnvFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	cfg := config.New()

	envPath, err := createProfile(path, cfg, tui.CreatedProfile{Name: "personal"})
	if err != nil {
		t.Fatalf("createProfile() error = %v", err)
	}
	if envPath != filepath.Join(dir, "personal.env") {
		t.Fatalf("envPath = %q, want personal.env path", envPath)
	}

	loaded, err := config.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile() error = %v", err)
	}
	profile := loaded.Profiles["personal"]
	if got, want := strings.Join(profile.EnvFiles, ","), "personal.env"; got != want {
		t.Fatalf("EnvFiles = %q, want %q", got, want)
	}

	envBytes, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("ReadFile(personal.env) error = %v", err)
	}
	if len(envBytes) != 0 {
		t.Fatalf("env file = %q, want empty", string(envBytes))
	}
}

func TestCreateProfileDoesNotOverwriteExistingEnvFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(filepath.Join(dir, "personal.env"), []byte("EXISTING=1\n"), 0o600); err != nil {
		t.Fatalf("write existing env file: %v", err)
	}

	_, err := createProfile(path, config.New(), tui.CreatedProfile{Name: "personal"})
	if err == nil {
		t.Fatal("createProfile() error = nil, want existing env file error")
	}
}

func TestOpenEditorRequiresEditor(t *testing.T) {
	err := openEditor("", "test.env", nil, nil, nil)
	if err == nil {
		t.Fatal("openEditor() error = nil, want EDITOR error")
	}
	if !strings.Contains(err.Error(), "EDITOR") {
		t.Fatalf("openEditor() error = %v, want EDITOR error", err)
	}
}

func TestProfileEnvPathUsesFirstRelativeEnvFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "work.env")
	if err := os.WriteFile(path, []byte("TOKEN=value\n"), 0o600); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	got, err := profileEnvPath(filepath.Join(dir, "config.toml"), config.Profile{
		EnvFiles: []string{"work.env", "common.env"},
	})
	if err != nil {
		t.Fatalf("profileEnvPath() error = %v", err)
	}
	if got != path {
		t.Fatalf("profileEnvPath() = %q, want %q", got, path)
	}
}

func TestProfileEnvPathRequiresExistingFile(t *testing.T) {
	_, err := profileEnvPath(filepath.Join(t.TempDir(), "config.toml"), config.Profile{
		EnvFiles: []string{"missing.env"},
	})
	if err == nil {
		t.Fatal("profileEnvPath() error = nil, want missing file error")
	}
}
