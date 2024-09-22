package main

import (
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
	"Pomodoro done",
	"Short Break done",
	"Long Break done",
	"Select",
}

type model struct {
	secondsRemaining int
	state            int
	cursor           int
	stateLengths     []int
	pomodoroStreak   int
}

func initialModel() model {
	return model{
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
	sb.WriteString(red("Pomodoro Timer 🍅\n\n"))

	if m.state == StateSelect {
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
		sb.WriteString(fmt.Sprintf("%s: %02d:%02d remaining\n", stateNames[m.state], m.secondsRemaining/60, m.secondsRemaining%60))
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("%s (streak: %d)\a\n", stateNames[m.state], m.pomodoroStreak))
	sb.WriteString("Next state: " + whiteBold(stateNames[m.getNextState()]) + grey(" [Enter to proceed...]"))

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
				m.state = m.cursor
				slog.Info("Setting", "state", stateNames[m.state])
				m.secondsRemaining = m.stateLengths[m.state]
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

			slog.Info("Finished", "state", stateNames[m.state])
			m.state = StatePomodoroDone + m.state
			if m.state == StatePomodoroDone {
				m.pomodoroStreak++
			}
		}

		return m, nil
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "enter", " ":
			m.state = m.getNextState()
			slog.Info("Setting", "state", stateNames[m.state])
			m.secondsRemaining = m.stateLengths[m.state]
			return m, doTick()
		case "esc":
			slog.Info("Escape")
			m.state = StateSelect
			m.cursor = 0
		}
	}

	return m, nil
}
