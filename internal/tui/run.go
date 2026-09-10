// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package tui

import (
	"io"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

type Input struct {
	Config config.Config
	Stdin  io.Reader
	Stdout io.Writer
}

type Result struct {
	Profile  string
	Created  *CreatedProfile
	Edit     string
	Canceled bool
}

type CreatedProfile struct {
	Name string
}

func Run(input Input) (Result, error) {
	model := NewModel(input.Config)
	var options []tea.ProgramOption
	if input.Stdin != nil {
		options = append(options, tea.WithInput(input.Stdin))
	}
	if input.Stdout != nil {
		options = append(options, tea.WithOutput(input.Stdout))
	}

	final, err := tea.NewProgram(model, options...).Run()
	if err != nil {
		return Result{}, err
	}
	if model, ok := final.(Model); ok {
		return model.Result(), nil
	}
	return Result{Canceled: true}, nil
}
