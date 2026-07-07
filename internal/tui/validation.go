// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package tui

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/config"
)

var profileNamePattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func sortedProfiles(cfg config.Config) []string {
	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func validateProfileName(name string, cfg config.Config) error {
	if name == "" {
		return fmt.Errorf("profile name is required")
	}
	if !profileNamePattern.MatchString(name) {
		return fmt.Errorf("profile name may only contain letters, numbers, dashes, and underscores")
	}
	if _, ok := cfg.Profiles[name]; ok {
		return fmt.Errorf("profile %q already exists", name)
	}
	return nil
}
