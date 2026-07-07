// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/config"
	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/tui"
)

func loadConfigForMode(path string, allowMissing bool) (config.Config, error) {
	cfg, err := config.LoadFile(path)
	if err == nil {
		return cfg, nil
	}
	if allowMissing && os.IsNotExist(err) {
		return config.New(), nil
	}
	return config.Config{}, err
}

func createProfile(path string, cfg config.Config, created tui.CreatedProfile) (string, error) {
	if _, ok := cfg.Profiles[created.Name]; ok {
		return "", fmt.Errorf("profile %q already exists", created.Name)
	}

	envFile := created.Name + ".env"
	envPath := filepath.Join(filepath.Dir(path), envFile)
	if _, err := os.Stat(envPath); err == nil {
		return "", fmt.Errorf("%s already exists", envPath)
	} else if !os.IsNotExist(err) {
		return "", err
	}

	cfg.Profiles[created.Name] = config.Profile{
		EnvFiles: []string{envFile},
		Env:      make(map[string]string),
	}

	if err := config.SaveFile(path, cfg); err != nil {
		return "", err
	}
	if err := createEnvFile(envPath); err != nil {
		return "", err
	}
	return envPath, nil
}

func createEnvFile(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	return file.Close()
}
