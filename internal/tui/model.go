// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package tui

import (
	"strings"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenList screen = iota
	screenName
)

type Model struct {
	cfg      config.Config
	names    []string
	cursor   int
	screen   screen
	input    textinput.Model
	err      string
	result   Result
	finished bool
}

func NewModel(cfg config.Config) Model {
	input := textinput.New()
	input.CharLimit = 256
	input.Width = 48

	return Model{
		cfg:    cfg,
		names:  sortedProfiles(cfg),
		input:  input,
		result: Result{Canceled: true},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			m.result = Result{Canceled: true}
			m.finished = true
			return m, tea.Quit
		}
		switch m.screen {
		case screenList:
			return m.updateList(msg)
		case screenName:
			return m.updateName(msg)
		}
	}
	return m, nil
}

func (m Model) Result() Result {
	return m.result
}

func (m Model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.err = ""
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.names)-1 {
			m.cursor++
		}
	case "enter":
		if len(m.names) == 0 {
			m.err = "no profiles configured; press n to create one"
			return m, nil
		}
		m.result = Result{Profile: m.names[m.cursor]}
		m.finished = true
		return m, tea.Quit
	case "n":
		m.screen = screenName
		m.input.Reset()
		m.input.Placeholder = "personal"
		m.input.Focus()
		return m, textinput.Blink
	}
	return m, nil
}

func (m Model) updateName(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		name := strings.TrimSpace(m.input.Value())
		if err := validateProfileName(name, m.cfg); err != nil {
			m.err = err.Error()
			return m, nil
		}
		m.result = Result{
			Profile: name,
			Created: &CreatedProfile{
				Name: name,
			},
		}
		m.finished = true
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}
