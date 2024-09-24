package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {
	logFile := flag.String("logfile", "/dev/null", "File to log into")
	dbFile := flag.String("dbfile", "", "SQLite database file (defaults to \"${HOME}/.config/pomodoro/pomodoro.db\")")
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

	if *dbFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			slog.Error(err.Error())
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err = os.MkdirAll(homeDir+"/.config/pomodoro", 0750); err != nil {
			slog.Error(err.Error())
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		*dbFile = homeDir + "/.config/pomodoro/pomodoro.db"
	}

	db, err = sql.Open("sqlite", *dbFile)
	if err != nil {
		slog.Error(err.Error())
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()

	if err = InitDB(db); err != nil {
		slog.Error(err.Error())
		fmt.Fprintln(os.Stderr, err)
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
