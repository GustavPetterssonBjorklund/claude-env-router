// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package main

import (
	"os"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/app"
)

func main() {
	os.Exit(app.RunWithInput(os.Args, os.Stdin, os.Stdout, os.Stderr))
}
