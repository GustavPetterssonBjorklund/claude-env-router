// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/config"
	cerenv "github.com/GustavPetterssonBjorklund/claude-env-router/internal/env"
	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/runner"
)

const binaryName = "cer"

// Run executes the CLI and returns a process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
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

	cfgPath, err := config.ResolvePath(opts.configPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	cfg, err := config.LoadFile(cfgPath)
	if err != nil {
		fmt.Fprintf(stderr, "load config: %v\n", err)
		return 1
	}

	if opts.profile == "" {
		printProfiles(stdout, cfg)
		printUsage(stdout)
		return 2
	}

	profile, ok := cfg.Profiles[opts.profile]
	if !ok {
		fmt.Fprintf(stderr, "unknown profile %q\n", opts.profile)
		printProfiles(stderr, cfg)
		return 1
	}

	envMap := cerenv.FromEnviron(os.Environ())
	baseDir := filepath.Dir(cfgPath)
	if err := cerenv.ApplyProfile(envMap, baseDir, profile); err != nil {
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

type options struct {
	configPath string
	profile    string
	claudeArgs []string
	help       bool
}

func parseArgs(args []string) (options, error) {
	var opts options
	if len(args) == 0 {
		args = []string{binaryName}
	}

	for i := 1; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			opts.help = true
		case "-c", "--config":
			i++
			if i >= len(args) {
				return opts, fmt.Errorf("%s requires a path", arg)
			}
			opts.configPath = args[i]
		case "--":
			if i+1 < len(args) {
				opts.claudeArgs = append(opts.claudeArgs, args[i+1:]...)
			}
			return opts, nil
		default:
			if strings.HasPrefix(arg, "-") {
				return opts, fmt.Errorf("unknown option %s", arg)
			}
			if opts.profile == "" {
				opts.profile = arg
			} else {
				opts.claudeArgs = append(opts.claudeArgs, arg)
			}
		}
	}

	return opts, nil
}

func printUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [--config path] <profile> [-- claude args...]\n", binaryName)
}

func printProfiles(w io.Writer, cfg config.Config) {
	if len(cfg.Profiles) == 0 {
		fmt.Fprintln(w, "No profiles configured.")
		return
	}

	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Fprintln(w, "Profiles:")
	for _, name := range names {
		fmt.Fprintf(w, "  %s\n", name)
	}
}
