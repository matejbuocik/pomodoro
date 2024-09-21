package main

import (
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: true})
	slog.SetDefault(slog.New(h))
	slog.Info("Starting")

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("Exiting")
}
