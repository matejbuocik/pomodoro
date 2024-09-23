package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Pomodoro struct {
	Id    string
	Type  int
	Start time.Time
	End   time.Time
	Note  string
}

func InitDB(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS pomodoro (
	id   	TEXT PRIMARY KEY,
	type 	INTEGER NOT NULL,
	start 	TEXT NOT NULL,
	end		TEXT NOT NULL,
	note    TEXT NOT NULL
	);`)
	return err
}

func AddPomodoro(db *sql.DB, p *Pomodoro) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	f := "2006-01-02T15:04:05 -0700"
	_, err = db.Exec(`insert into pomodoro(id, type, start, end, note) values (?, ?, ?, ?, ?)`,
		uuid.String(), p.Type, p.Start.Format(f), p.End.Format(f), p.Note)

	return err
}
