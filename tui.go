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

type model struct {
	secondsRemaining int
	state            int
	choices          []string
	cursor           int
}

func (m model) View() string {
	var sb strings.Builder
	sb.WriteString(red("Pomodoro Timer 🍅\n\n"))

	switch m.state {
	case 0:
		for i, choice := range m.choices {
			if m.cursor == i {
				sb.WriteString(whiteBold(fmt.Sprintf("> %s\n", choice)))
			} else {
				sb.WriteString(white(fmt.Sprintf("  %s\n", choice)))
			}
		}

	case 1, 2, 3:
		sb.WriteString(fmt.Sprintf("%02d:%02d remaining\n", m.secondsRemaining/60, m.secondsRemaining%60))
	}

	return sb.String()
}

func initialModel() model {
	return model{
		choices: []string{
			"Pomodoro    (25min)",
			"Short Break (5min)",
			"Long Break  (15min)",
		},
	}
}

type TickMsg struct{}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Pomodoro")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			slog.Info("Setting", "state", m.cursor+1)

			switch m.cursor {
			case 0:
				m.state = 1
				m.secondsRemaining = 25 * 60
			case 1:
				m.state = 2
				m.secondsRemaining = 5 * 60
			case 2:
				m.state = 3
				m.secondsRemaining = 15 * 60
			}

			return m, doTick()

		case "esc":
			if m.state != 0 {
				slog.Info("Escape")
				m.state = 0
				m.cursor = 0
			}
		}

	case TickMsg:
		if m.state == 0 {
			break
		}

		m.secondsRemaining--
		if m.secondsRemaining <= 0 {
			slog.Info("Finished", "state", m.state)
			m.state = 0
			m.cursor = 0
		} else {
			return m, doTick()
		}
	}

	return m, nil
}
