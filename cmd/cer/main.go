// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package main

import (
	"os"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args, os.Stdout, os.Stderr))
}
