// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package app

import (
	"fmt"
	"io"
	"sort"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/config"
)

func printUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [--config path] <profile> [-- claude args...]\n", environment.binaryName)
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
