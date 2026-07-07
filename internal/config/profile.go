// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package config

type Profile struct {
	EnvFiles []string
	Env      map[string]string
	Args     []string
}

func NewProfile() Profile {
	return Profile{
		Env: make(map[string]string),
	}
}
