// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package env

import "os"

func Expand(value string, env map[string]string) string {
	return os.Expand(value, func(key string) string {
		return env[key]
	})
}
