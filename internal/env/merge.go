// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package env

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/config"
	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/vault"
)

func FromEnviron(environ []string) map[string]string {
	out := make(map[string]string, len(environ))
	for _, entry := range environ {
		key, value, ok := strings.Cut(entry, "=")
		if ok {
			out[key] = value
		}
	}
	return out
}

func ToEnviron(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for key := range env {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	out := make([]string, 0, len(keys))
	for _, key := range keys {
		out = append(out, key+"="+env[key])
	}
	return out
}

func ApplyProfile(target map[string]string, baseDir string, profile config.Profile, secrets map[string]string) error {
	for _, envFile := range profile.EnvFiles {
		path := envFile
		if !filepath.IsAbs(path) {
			path = filepath.Join(baseDir, path)
		}
		values, err := LoadDotenv(path)
		if err != nil {
			return err
		}
		for key, value := range values {
			target[key] = Expand(value, target)
		}
	}

	for key, value := range profile.Env {
		target[key] = Expand(value, target)
	}
	for key, value := range secrets {
		target[key] = Expand(value, target)
	}

	return nil
}

func LoadSecrets(path, profile string) (map[string]string, error) {
	store := vault.Store{Path: path}
	values, err := store.List(profile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	out := make(map[string]string, len(values))
	for _, key := range values {
		value, err := store.Get(profile, key)
		if err != nil {
			return nil, err
		}
		out[key] = value
	}
	return out, nil
}

func LoadDotenv(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	defer file.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("%s:%d: expected KEY=value", path, lineNo)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			return nil, fmt.Errorf("%s:%d: empty key", path, lineNo)
		}
		values[key] = trimValue(value)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return values, nil
}

func trimValue(value string) string {
	if len(value) >= 2 {
		first := value[0]
		last := value[len(value)-1]
		if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
			return value[1 : len(value)-1]
		}
	}
	return value
}
