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
}

func parseArgs(args []string) (options, error) {
	var opts options
	if len(args) == 0 {
		args = []string{environment.binaryName}
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
