// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package detect

import (
	"os"
	"path/filepath"
)

func FindUp(name, start string) (string, bool) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", false
	}

	for {
		candidate := filepath.Join(dir, name)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}
