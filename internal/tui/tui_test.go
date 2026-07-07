// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package tui

import (
	"testing"

	"github.com/GustavPetterssonBjorklund/claude-env-router/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

func TestModelSelectsProfile(t *testing.T) {
	cfg := config.New()
	cfg.Profiles["work"] = config.NewProfile()
	cfg.Profiles["personal"] = config.NewProfile()

	model := NewModel(cfg)
	model = updateModel(t, model, key("down"))
	model = updateModel(t, model, key("enter"))

	result := model.Result()
	if result.Profile != "work" {
		t.Fatalf("selected profile = %q, want work", result.Profile)
	}
	if result.Canceled {
		t.Fatal("result is canceled")
	}
}

func TestModelCreatesProfile(t *testing.T) {
	model := NewModel(config.New())
	model = updateModel(t, model, key("n"))

	model.input.SetValue("personal")
	model = updateModel(t, model, key("enter"))

	result := model.Result()
	if result.Profile != "personal" {
		t.Fatalf("profile = %q, want personal", result.Profile)
	}
	if result.Created == nil {
		t.Fatal("Created is nil")
	}
	if result.Created.Name != "personal" {
		t.Fatalf("Created.Name = %q, want personal", result.Created.Name)
	}
	if result.Canceled {
		t.Fatal("result is canceled")
	}
}

func TestModelRejectsDuplicateProfile(t *testing.T) {
	cfg := config.New()
	cfg.Profiles["personal"] = config.NewProfile()

	model := NewModel(cfg)
	model = updateModel(t, model, key("n"))
	model.input.SetValue("personal")
	model = updateModel(t, model, key("enter"))

	if model.screen != screenName {
		t.Fatalf("screen = %v, want screenName", model.screen)
	}
	if model.err == "" {
		t.Fatal("expected duplicate profile error")
	}
}

func updateModel(t *testing.T, model Model, msg tea.Msg) Model {
	t.Helper()
	updated, _ := model.Update(msg)
	next, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model has type %T, want Model", updated)
	}
	return next
}

func key(value string) tea.KeyMsg {
	switch value {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(value)}
	}
}
