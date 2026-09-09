// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package app

import (
	"fmt"
	"strings"
)

type options struct {
	configPath string
	profile    string
	claudeArgs []string
	help       bool
	secretOp   string
	secretKey  string
}

func parseArgs(args []string) (options, error) {
	var opts options
	if len(args) == 0 {
		args = []string{environment.binaryName}
	}

	for i := 1; i < len(args); i++ {
		arg := args[i]
		if arg == "secret" || arg == "secrets" {
			if i+1 >= len(args) {
				return opts, fmt.Errorf("%s requires an operation", arg)
			}
			op := args[i+1]
			if op == "migrate" {
				opts.secretOp = op
				return opts, nil
			}
			if op != "set" && op != "unset" && op != "list" {
				return opts, fmt.Errorf("unknown secret operation %s", op)
			}
			if i+2 >= len(args) {
				return opts, fmt.Errorf("secret %s requires a profile", op)
			}
			opts.secretOp = op
			opts.profile = args[i+2]
			if op != "list" {
				if i+3 >= len(args) {
					return opts, fmt.Errorf("secret %s requires a key", op)
				}
				opts.secretKey = args[i+3]
			}
			return opts, nil
		}
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
