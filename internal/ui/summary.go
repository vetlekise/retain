package ui

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/vetlekise/retain/internal/session"
)

type SummaryModel struct {
	session      *session.Session // pointer — shared, mutable
	title        string           // exam title from meta.title
	elapsed      time.Duration    // time snapshot captured at transition, not live
	scrollOffset int              // index of the first visible question row
	width        int              // terminal width (from WindowSizeMsg)
	height       int              // terminal height (from WindowSizeMsg)
}

func NewSummaryModel(s *session.Session, title string) SummaryModel {
	return SummaryModel{
		session: s,
		title:   title,
		// snapshot elapsed here so the summary shows the time at completion,
		// not a live-ticking value that keeps growing while the user reads it.
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
			// -1 so the last question is at most the first visible row, not scrolled off screen entirely
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

	// fixedLines is the number of lines in the view that are not part of the
	// scrollable question list: 2 title/blank, 1 score, 1 time, 1 blank (5
	// header lines) + 1 top padding + 1 bottom padding + 1 footer = 8.
	// Subtract from terminal height to get how many question rows fit on screen.
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

	// m.width is 0 until the first WindowSizeMsg arrives; guard against
	// rendering a zero-width styled box on the very first frame.
	if m.width > 0 {
		v := tea.NewView(bodyStyle(m.width, m.height).Render(body) + "\n" + footerStyle(m.width).Render(footer))
		v.AltScreen = true
		return v
	}

	// pre-WindowSizeMsg fallback: plain output, no styling or AltScreen.
	return tea.NewView(body + "\n" + footer + "\n")
}
