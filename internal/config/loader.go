// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/detect"
)

func ResolvePath(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	if envPath := os.Getenv(EnvConfigPath); envPath != "" {
		return envPath, nil
	}
	if found, ok := detect.FindUp("cer.toml", "."); ok {
		return found, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("find home directory: %w", err)
	}
	return filepath.Join(home, ".config", "cer", "config.toml"), nil
}

func LoadFile(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	cfg, err := Parse(file)
	if err != nil {
		return Config{}, fmt.Errorf("%s: %w", path, err)
	}
	return cfg, nil
}

func SaveFile(path string, cfg Config) error {
	var buf bytes.Buffer
	Write(&buf, cfg)

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o600)
}

func Write(w io.Writer, cfg Config) {
	fmt.Fprintf(w, "binary = %s\n", formatString(cfg.Binary))

	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		profile := cfg.Profiles[name]
		fmt.Fprintf(w, "\n[profiles.%s]\n", name)
		fmt.Fprintf(w, "env_files = %s\n", formatStringArray(profile.EnvFiles))
		fmt.Fprintf(w, "args = %s\n", formatStringArray(profile.Args))

		if len(profile.Env) == 0 {
			continue
		}

		fmt.Fprintf(w, "\n[profiles.%s.env]\n", name)
		keys := make([]string, 0, len(profile.Env))
		for key := range profile.Env {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			fmt.Fprintf(w, "%s = %s\n", key, formatString(profile.Env[key]))
		}
	}
}

func Parse(file io.Reader) (Config, error) {
	cfg := New()
	scanner := bufio.NewScanner(file)
	section := ""
	lineNo := 0

	for scanner.Scan() {
		lineNo++
		line := stripComment(strings.TrimSpace(scanner.Text()))
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"))
			if strings.HasPrefix(section, "profiles.") {
				name := profileName(section)
				if name != "" {
					ensureProfile(cfg, name)
				}
			}
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return Config{}, fmt.Errorf("line %d: expected key = value", lineNo)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		if section == "" {
			if key == "binary" {
				parsed, err := parseString(value)
				if err != nil {
					return Config{}, fmt.Errorf("line %d: %w", lineNo, err)
				}
				cfg.Binary = parsed
			}
			continue
		}

		if !strings.HasPrefix(section, "profiles.") {
			continue
		}

		name := profileName(section)
		if name == "" {
			return Config{}, fmt.Errorf("line %d: invalid profile section %q", lineNo, section)
		}
		profile := ensureProfile(cfg, name)

		switch {
		case strings.HasSuffix(section, ".env"):
			parsed, err := parseString(value)
			if err != nil {
				return Config{}, fmt.Errorf("line %d: %w", lineNo, err)
			}
			profile.Env[key] = parsed
		case key == "env_files":
			parsed, err := parseStringArray(value)
			if err != nil {
				return Config{}, fmt.Errorf("line %d: %w", lineNo, err)
			}
			profile.EnvFiles = parsed
		case key == "args":
			parsed, err := parseStringArray(value)
			if err != nil {
				return Config{}, fmt.Errorf("line %d: %w", lineNo, err)
			}
			profile.Args = parsed
		}

		cfg.Profiles[name] = profile
	}

	if err := scanner.Err(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func ensureProfile(cfg Config, name string) Profile {
	profile, ok := cfg.Profiles[name]
	if !ok {
		profile = NewProfile()
		cfg.Profiles[name] = profile
	}
	if profile.Env == nil {
		profile.Env = make(map[string]string)
	}
	return profile
}

func profileName(section string) string {
	rest := strings.TrimPrefix(section, "profiles.")
	name, _, _ := strings.Cut(rest, ".")
	return name
}

func stripComment(line string) string {
	inQuote := false
	escaped := false
	for i, r := range line {
		if escaped {
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if r == '"' {
			inQuote = !inQuote
			continue
		}
		if r == '#' && !inQuote {
			return strings.TrimSpace(line[:i])
		}
	}
	return line
}

func parseString(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		parsed := strings.ReplaceAll(value[1:len(value)-1], `\"`, `"`)
		parsed = strings.ReplaceAll(parsed, `\\`, `\`)
		return parsed, nil
	}
	if value == "" {
		return "", fmt.Errorf("empty string value")
	}
	return value, nil
}

func parseStringArray(value string) ([]string, error) {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "[") || !strings.HasSuffix(value, "]") {
		return nil, fmt.Errorf("expected string array")
	}
	value = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(value, "["), "]"))
	if value == "" {
		return nil, nil
	}

	var out []string
	for _, part := range splitArrayParts(value) {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, fmt.Errorf("empty array element")
		}
		parsed, err := parseString(part)
		if err != nil {
			return nil, err
		}
		out = append(out, parsed)
	}
	return out, nil
}

func splitArrayParts(value string) []string {
	var parts []string
	start := 0
	inQuote := false
	escaped := false
	for i, r := range value {
		if escaped {
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if r == '"' {
			inQuote = !inQuote
			continue
		}
		if r == ',' && !inQuote {
			parts = append(parts, value[start:i])
			start = i + 1
		}
	}
	return append(parts, value[start:])
}

func formatString(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `"`, `\"`)
	return `"` + value + `"`
}

func formatStringArray(values []string) string {
	if len(values) == 0 {
		return "[]"
	}
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, formatString(value))
	}
	return "[" + strings.Join(parts, ", ") + "]"
}
