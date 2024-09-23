package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func red(s string) string {
	return fmt.Sprintf("\033[31;1m%s\033[0m", s)
}

func white(s string) string {
	return fmt.Sprintf("\033[37m%s\033[0m", s)
}

func whiteBold(s string) string {
	return fmt.Sprintf("\033[37;1m%s\033[0m", s)
}

func grey(s string) string {
	return fmt.Sprintf("\033[30;1m%s\033[0m", s)
}

func green(s string) string {
	return fmt.Sprintf("\033[32;1m%s\033[0m", s)
}

func cyan(s string) string {
	return fmt.Sprintf("\033[36;1m%s\033[0m", s)
}

const (
	StatePomodoro = iota
	StateShortBreak
	StateLongBreak

	StatePomodoroDone
	StateShortBreakDone
	StateLongBreakDone

	StateSelect
)

var stateNames = []string{
	"Pomodoro",
	"Short Break",
	"Long Break",
}

type model struct {
	secondsRemaining int
	state            int
	cursor           int
	stateLengths     []int
	pomodoroStreak   int

	db              *sql.DB
	currentPomodoro *Pomodoro
}

func initialModel(db *sql.DB) model {
	return model{
		db:    db,
		state: StateSelect,
		stateLengths: []int{
			25 * 60,
			5 * 60,
			15 * 60,
		},
	}
}

func (m model) getNextState() int {
	next := StateSelect

	switch m.state {
	case StatePomodoro:
		next = StatePomodoroDone
	case StateShortBreak:
		next = StateShortBreak
	case StateLongBreak:
		next = StateLongBreakDone
	case StatePomodoroDone:
		if m.pomodoroStreak > 0 && m.pomodoroStreak%4 == 0 {
			next = StateLongBreak
		} else {
			next = StateShortBreak
		}
	case StateShortBreakDone:
		next = StatePomodoro
	case StateLongBreakDone:
		next = StatePomodoro
	}

	return next
}

func (m model) View() string {
	var sb strings.Builder

	if m.state == StateSelect {
		sb.WriteString(red("Pomodoro Timer 🍅\n\n"))
		for i, choice := range stateNames[StatePomodoro : StateLongBreak+1] {
			if m.cursor == i {
				sb.WriteString(whiteBold(fmt.Sprintf("> %s\t(%d min)\n", choice, m.stateLengths[i]/60)))
			} else {
				sb.WriteString(white(fmt.Sprintf("  %s\t(%d min)\n", choice, m.stateLengths[i]/60)))
			}
		}
		return sb.String()
	}

	if m.state <= StateLongBreak {
		msg := "Focus!"
		if m.state > StatePomodoro {
			msg = "Chill."
		}
		sb.WriteString(fmt.Sprintf("%s: %02d:%02d remaining\n", whiteBold(msg), m.secondsRemaining/60, m.secondsRemaining%60))
		return sb.String()
	}

	if m.state == StatePomodoroDone {
		sb.WriteString(green("Pomodoro complete "))
		sb.WriteString(fmt.Sprintf("✅\nGood job! (streak: %d)\a\n", m.pomodoroStreak))
	} else {
		sb.WriteString(cyan("Chilled out and ready for the next one!\a\n"))
	}
	sb.WriteString("Upcoming: " + whiteBold(stateNames[m.getNextState()]) + grey(" [Enter to proceed...]"))

	return sb.String()
}

type TickMsg struct{}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Pomodoro 🍅")
}

func (m *model) startPomodoro(t int) {
	m.state = t
	slog.Info("Start", "state", stateNames[m.state])
	m.secondsRemaining = m.stateLengths[m.state]
	m.currentPomodoro = &Pomodoro{
		Type:  m.state,
		Start: time.Now(),
	}
}

func (m *model) endPomodoro() {
	slog.Info("End", "state", stateNames[m.state])
	m.currentPomodoro.End = time.Now()

	go func(db *sql.DB, p *Pomodoro, stateName string) {
		if err := AddPomodoro(db, p); err != nil {
			slog.Error(err.Error())
		} else {
			slog.Info("DB Save", "state", stateName)
		}
	}(m.db, m.currentPomodoro, stateNames[m.state])

	m.currentPomodoro = nil
	m.state = StatePomodoroDone + m.state
	if m.state == StatePomodoroDone {
		m.pomodoroStreak++
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	if m.state == StateSelect {
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				if m.cursor < 2 {
					m.cursor++
				}
			case "enter", " ":
				m.startPomodoro(m.cursor)
				return m, doTick()
			}
		}

		return m, nil
	}

	if StatePomodoro <= m.state && m.state <= StateLongBreak {
		if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "esc" {
			slog.Info("Escape")
			m.state = StateSelect
			m.cursor = 0
		}

		if _, ok := msg.(TickMsg); ok {
			m.secondsRemaining--
			if m.secondsRemaining > 0 {
				return m, doTick()
			}

			m.endPomodoro()
		}

		return m, nil
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "enter", " ":
			m.startPomodoro(m.getNextState())
			return m, doTick()
		case "esc":
			slog.Info("Escape")
			m.state = StateSelect
			m.cursor = 0
		}
	}

	return m, nil
}
