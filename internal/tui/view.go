// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
	mutedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

func (m Model) View() string {
	if m.finished {
		return ""
	}

	var b strings.Builder
	fmt.Fprintln(&b, titleStyle.Render("cer"))
	fmt.Fprintln(&b)

	switch m.screen {
	case screenList:
		m.viewList(&b)
	case screenName:
		m.viewPrompt(&b, "New profile name")
	}

	if m.err != "" {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, errorStyle.Render(m.err))
	}
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, mutedStyle.Render("Enter select/confirm  e edit env  n new env  Esc quit"))
	return b.String()
}

func (m Model) viewList(b *strings.Builder) {
	if len(m.names) == 0 {
		fmt.Fprintln(b, mutedStyle.Render("No profiles configured."))
		fmt.Fprintln(b)
		fmt.Fprintln(b, "Press n to create a new env.")
		return
	}

	fmt.Fprintln(b, "Choose a profile:")
	fmt.Fprintln(b)
	for i, name := range m.names {
		line := "  " + name
		if i == m.cursor {
			line = selectedStyle.Render("> " + name)
		}
		fmt.Fprintln(b, line)
	}
}

func (m Model) viewPrompt(b *strings.Builder, label string) {
	fmt.Fprintln(b, label)
	fmt.Fprintln(b, m.input.View())
}
