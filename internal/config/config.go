// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package config

const DefaultClaudeBinary = "claude"

type Config struct {
	Binary   string
	Profiles map[string]Profile
}

func New() Config {
	return Config{
		Binary:   DefaultClaudeBinary,
		Profiles: make(map[string]Profile),
	}
}
