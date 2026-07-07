// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/detect"
)

const EnvConfigPath = "CER_CONFIG"

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
		return strings.ReplaceAll(value[1:len(value)-1], `\"`, `"`), nil
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
	for _, part := range strings.Split(value, ",") {
		parsed, err := parseString(strings.TrimSpace(part))
		if err != nil {
			return nil, err
		}
		out = append(out, parsed)
	}
	return out, nil
}
