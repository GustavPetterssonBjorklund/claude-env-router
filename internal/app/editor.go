// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func openEditor(editor, path string, stdin io.Reader, stdout, stderr io.Writer) error {
	fields := strings.Fields(editor)
	if len(fields) == 0 {
		return fmt.Errorf("EDITOR is not set")
	}

	cmd := exec.Command(fields[0], append(fields[1:], path)...)
	if file, ok := stdin.(*os.File); ok {
		cmd.Stdin = file
	}
	if file, ok := stdout.(*os.File); ok {
		cmd.Stdout = file
	}
	if file, ok := stderr.(*os.File); ok {
		cmd.Stderr = file
	}
	return cmd.Run()
}

func editorConfigured(editor string) bool {
	return len(strings.Fields(editor)) > 0
}
