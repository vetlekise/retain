package ui

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/vetlekise/retain/internal/session"
)

type SummaryModel struct {
	session      *session.Session
	title        string
	elapsed      time.Duration
	scrollOffset int
	width        int
	height       int
}

func NewSummaryModel(s *session.Session, title string) SummaryModel {
	return SummaryModel{
		session: s,
		title:   title,
		elapsed: s.Elapsed(),
	}
}

func (m SummaryModel) Init() tea.Cmd {
	return nil
}

func (m SummaryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			m.session.Reset()
			return NewExamModel(m.session, m.title), nil
		case "down", "j":
			max := len(m.session.Questions) - 1
			if m.scrollOffset < max {
				m.scrollOffset++
			}
		case "up", "k":
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}
		}
	}

	return m, nil
}

func (m SummaryModel) View() tea.View {
	correct, total := m.session.Score()
	pct := 0
	if total > 0 {
		pct = correct * 100 / total
	}

	header := titleStyle.Render("Exam Complete") + "\n\n"
	header += fmt.Sprintf("Score: %d/%d (%d%%)\n", correct, total, pct)
	header += fmt.Sprintf("Time:  %s\n\n", m.elapsed.Round(time.Second))

	footer := "j/k: scroll · r: retry · q: quit"

	// 5 header lines + 2 padding (top+bottom) + 1 footer line
	const fixedLines = 8
	available := m.height - fixedLines
	if available < 1 {
		available = 1
	}

	questions := m.session.Questions
	start := m.scrollOffset
	end := start + available
	if end > len(questions) {
		end = len(questions)
	}

	rows := ""
	for i := start; i < end; i++ {
		q := questions[i]
		mark := correctStyle.Render("✓")
		if !m.session.IsCorrect(q.ID) {
			mark = incorrectStyle.Render("✗")
		}
		rows += fmt.Sprintf("%s  Q%d: %s\n", mark, i+1, q.Text)
	}

	if end < len(questions) {
		rows += "↓ more\n"
	}

	body := header + rows

	if m.width > 0 {
		v := tea.NewView(bodyStyle(m.width, m.height).Render(body) + "\n" + footerStyle(m.width).Render(footer))
		v.AltScreen = true
		return v
	}

	return tea.NewView(body + "\n" + footer + "\n")
}
