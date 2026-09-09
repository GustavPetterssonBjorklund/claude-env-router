// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/config"
	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/vault"
)

func runSecretCommand(opts options, cfg config.Config, stdin io.Reader, stdout, stderr io.Writer) int {
	if opts.secretOp == "migrate" {
		return migrateSecrets(cfg, stderr)
	}
	if _, ok := cfg.Profiles[opts.profile]; !ok {
		fmt.Fprintf(stderr, "unknown profile %q\n", opts.profile)
		return 1
	}
	store := vault.Store{Path: cfg.FilePaths.VaultPath}
	switch opts.secretOp {
	case "list":
		names, err := store.List(opts.profile)
		if errors.Is(err, os.ErrNotExist) {
			return 0
		}
		if err != nil {
			fmt.Fprintf(stderr, "list secrets: %v\n", err)
			return 1
		}
		for _, name := range names {
			fmt.Fprintln(stdout, name)
		}
		return 0
	case "migrate":
		return migrateSecrets(cfg, stderr)
	case "set":
		value, err := readSecret(stdin, stderr, opts.secretKey)
		if err != nil {
			fmt.Fprintf(stderr, "read secret: %v\n", err)
			return 1
		}
		if err := store.Set(opts.profile, opts.secretKey, value); err != nil {
			fmt.Fprintf(stderr, "save secret: %v\n", err)
			return 1
		}
		return 0
	case "unset":
		if err := store.Delete(opts.profile, opts.secretKey); err != nil {
			if errors.Is(err, vault.ErrNotFound) || errors.Is(err, os.ErrNotExist) {
				fmt.Fprintf(stderr, "secret %q does not exist\n", opts.secretKey)
				return 1
			}
			fmt.Fprintf(stderr, "delete secret: %v\n", err)
			return 1
		}
		return 0
	default:
		return 2
	}
}

func migrateSecrets(cfg config.Config, stderr io.Writer) int {
	store := vault.Store{Path: cfg.FilePaths.VaultPath}
	for profileName, profile := range cfg.Profiles {
		if len(profile.Env) == 0 {
			continue
		}
		for key, value := range profile.Env {
			if err := store.Set(profileName, key, value); err != nil {
				fmt.Fprintf(stderr, "migrate %s.%s: %v\n", profileName, key, err)
				return 1
			}
		}
		profile.Env = nil
		cfg.Profiles[profileName] = profile
	}
	if err := config.SaveFile(cfg.FilePaths.ConfigPath, cfg); err != nil {
		fmt.Fprintf(stderr, "save migrated config: %v\n", err)
		return 1
	}
	return 0
}

func readSecret(stdin io.Reader, stderr io.Writer, key string) (string, error) {
	fmt.Fprintf(stderr, "Value for %s: ", key)
	if stdin == os.Stdin {
		value, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(stderr)
		return string(value), err
	}
	reader := bufio.NewReader(stdin)
	value, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return strings.TrimSuffix(strings.TrimSuffix(value, "\n"), "\r"), nil
}
