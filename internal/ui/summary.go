package ui

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/vetlekise/retain/internal/session"
)

type SummaryModel struct {
	session *session.Session
	title   string
	elapsed time.Duration
	width   int
	height  int
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

	body := titleStyle.Render("Exam Complete") + "\n\n"
	body += fmt.Sprintf("Score: %d/%d (%d%%)\n", correct, total, pct)
	body += fmt.Sprintf("Time:  %s\n\n", m.elapsed.Round(time.Second))

	for i, q := range m.session.Questions {
		mark := correctStyle.Render("✓")
		if !m.session.IsCorrect(q.ID) {
			mark = incorrectStyle.Render("✗")
		}
		body += fmt.Sprintf("%s  Q%d: %s\n", mark, i+1, q.Text)
	}

	footer := "r: retry · q: quit"

	if m.width > 0 {
		v := tea.NewView(bodyStyle(m.width, m.height).Render(body) + "\n" + footerStyle(m.width).Render(footer))
		v.AltScreen = true
		return v
	}

	return tea.NewView(body + "\n" + footer + "\n")
}
