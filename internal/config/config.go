// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultClaudeBinaryName = "claude"
	DefaultRouterBinaryName = "cer"
	DefaultProfileName      = "default"
	DefaultProfilesDir      = "profiles"
	DefaultConfigDirName    = "cer"
	DefaultConfigFileName   = "config.toml"
	EnvConfigPath           = "CER_CONFIG"
)

var (
	ClaudeBinary          = envOrDefault("CLAUDE_BINARY", DefaultClaudeBinaryName)
	ClaudeEnvRouterBinary = envOrDefault(
		"CLAUDE_ENV_ROUTER_BINARY",
		DefaultRouterBinaryName,
	)
	ProfileName = envOrDefault(
		"CLAUDE_ENV_ROUTER_DEFAULT_PROFILE",
		DefaultProfileName,
	)
	ProfilesDir = envOrDefault(
		"CLAUDE_ENV_ROUTER_PROFILES_DIR",
		DefaultProfilesDir,
	)
	ConfigPath = envOrDefault(
		EnvConfigPath,
		defaultConfigPath(),
	)
)

func envOrDefault(name, fallback string) string {
	value, exists := os.LookupEnv(name)
	if !exists || value == "" {
		return fallback
	}

	return value
}

func defaultConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		homeDir, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return filepath.Join(".config", DefaultConfigDirName, DefaultConfigFileName)
		}

		return filepath.Join(homeDir, ".config", DefaultConfigDirName, DefaultConfigFileName)
	}

	return filepath.Join(configDir, DefaultConfigDirName, DefaultConfigFileName)
}

type FilePaths struct {
	ConfigPath  string
	ProfilesDir string
}

type Config struct {
	ClaudeEnvRouterBinary string
	Binary                string
	FilePaths             FilePaths
	Profiles              map[string]Profile
}

func New() Config {
	return Config{
		ClaudeEnvRouterBinary: ClaudeEnvRouterBinary,
		Binary:                ClaudeBinary,
		FilePaths: FilePaths{
			ConfigPath:  ConfigPath,
			ProfilesDir: ProfilesDir,
		},
		Profiles: make(map[string]Profile),
	}
}
