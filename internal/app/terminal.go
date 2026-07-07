// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package app

import (
	"io"
	"os"
)

func isInteractive(stdin io.Reader, stdout io.Writer) bool {
	in, ok := stdin.(*os.File)
	if !ok {
		return false
	}
	out, ok := stdout.(*os.File)
	if !ok {
		return false
	}
	return isTerminal(in) && isTerminal(out)
}

func isTerminal(file *os.File) bool {
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
