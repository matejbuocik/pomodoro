package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {
	logFile := flag.String("logfile", "/dev/null", "File to log into")
	dbFile := flag.String("dbfile", "pomodoro.db", "Database file")
	flag.Parse()

	logF, err := os.Create(*logFile)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer logF.Close()

	h := slog.NewTextHandler(logF, &slog.HandlerOptions{AddSource: true})
	slog.SetDefault(slog.New(h))
	slog.Info("Init start")

	db, err = sql.Open("sqlite", *dbFile)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	if err = InitDB(db); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("Init done")

	p := tea.NewProgram(initialModel(db))
	if _, err := p.Run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("Exit")
}
