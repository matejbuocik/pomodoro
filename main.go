package main

import (
	"flag"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logFile := flag.String("logfile", "/dev/null", "File to log into")
	flag.Parse()

	file, err := os.Create(*logFile)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer file.Close()

	h := slog.NewTextHandler(file, &slog.HandlerOptions{AddSource: true})
	slog.SetDefault(slog.New(h))
	slog.Info("Starting")

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("Exiting")
}
