// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/config"
	cerenv "github.com/GustavPetterssonBjorklund/claude-env-router/internal/env"
	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/runner"
	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/tui"
)

type Environment struct {
	binaryName     string
	claudeProcName string
}

var environment = Environment{
	binaryName:     "cer",
	claudeProcName: "claude",
}

// Run executes the CLI and returns a process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	return RunWithInput(args, os.Stdin, stdout, stderr)
}

func RunWithInput(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	opts, err := parseArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		printUsage(stderr)
		return 2
	}

	if opts.help {
		printUsage(stdout)
		return 0
	}

	if opts.profile == "" && !isInteractive(stdin, stdout) {
		fmt.Fprintln(stderr, "no profile specified; pass a profile or run cer interactively")
		return 2
	}

	cfgPath, err := config.ResolvePath(opts.configPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	cfg, err := loadConfigForMode(cfgPath, opts.profile == "")
	if err != nil {
		fmt.Fprintf(stderr, "load config: %v\n", err)
		return 1
	}
	cfg.FilePaths.ConfigPath = cfgPath
	cfg.FilePaths.VaultPath = filepath.Join(filepath.Dir(cfgPath), config.DefaultVaultFileName)

	if opts.secretOp != "" {
		return runSecretCommand(opts, cfg, stdin, stdout, stderr)
	}

	if opts.profile == "" {
		result, err := tui.Run(tui.Input{
			Config: cfg,
			Stdin:  stdin,
			Stdout: stdout,
		})
		if err != nil {
			fmt.Fprintf(stderr, "run tui: %v\n", err)
			return 1
		}
		if result.Canceled {
			return 0
		}
		if result.Created != nil {
			editor := os.Getenv("EDITOR")
			if !editorConfigured(editor) {
				fmt.Fprintln(stderr, "EDITOR is not set")
				return 1
			}
			envPath, err := createProfile(cfgPath, cfg, *result.Created)
			if err != nil {
				fmt.Fprintf(stderr, "create profile: %v\n", err)
				return 1
			}
			if err := openEditor(editor, envPath, stdin, stdout, stderr); err != nil {
				fmt.Fprintf(stderr, "edit %s: %v\n", envPath, err)
				return 1
			}
			cfg, err = config.LoadFile(cfgPath)
			if err != nil {
				fmt.Fprintf(stderr, "load config: %v\n", err)
				return 1
			}
		}
		opts.profile = result.Profile
	}

	profile, ok := cfg.Profiles[opts.profile]
	if !ok {
		fmt.Fprintf(stderr, "unknown profile %q\n", opts.profile)
		printProfiles(stderr, cfg)
		return 1
	}

	envMap := cerenv.FromEnviron(os.Environ())
	baseDir := filepath.Dir(cfgPath)
	secrets, err := cerenv.LoadSecrets(cfg.FilePaths.VaultPath, opts.profile)
	if err != nil {
		fmt.Fprintf(stderr, "load secrets: %v\n", err)
		return 1
	}
	if err := cerenv.ApplyProfile(envMap, baseDir, profile, secrets); err != nil {
		fmt.Fprintf(stderr, "prepare environment: %v\n", err)
		return 1
	}

	claudeArgs := append([]string{}, profile.Args...)
	claudeArgs = append(claudeArgs, opts.claudeArgs...)

	err = runner.Exec(context.Background(), cfg.Binary, claudeArgs, cerenv.ToEnviron(envMap), "", stdout, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "run %s: %v\n", cfg.Binary, err)
		return 1
	}

	return 0
}
