package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/vetlekise/retain/internal/session"
)

type ExamModel struct {
	session   *session.Session // pointer — shared, mutable
	title     string           // exam title from meta.title
	submitted bool             // has the current question been submitted?
	pending   map[string]bool  // option ID → selected
	cursor    int              // which option the > marker is on
	showHint  bool             // is the hint visible?
	width     int              // terminal width (from WindowSizeMsg)
	height    int              // terminal height (from WindowSizeMsg)
}

func NewExamModel(s *session.Session, title string) ExamModel {
	return ExamModel{
		session: s,
		title:   title,
		pending: make(map[string]bool),
	}
}

func (m ExamModel) Init() tea.Cmd {
	return nil
}

func (m ExamModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.session.Current().Options)-1 {
				m.cursor++
			}

		case "space":
			optID := m.session.Current().Options[m.cursor].ID
			m.pending[optID] = !m.pending[optID]

		case "enter":
			if !m.submitted {
				q := m.session.Current()
				var selected []string
				for _, opt := range q.Options {
					if m.pending[opt.ID] {
						selected = append(selected, opt.ID)
					}
				}
				m.session.Answer(q.ID, selected)
				m.submitted = true
				m.showHint = false
			} else {
				if m.session.CurrentIndex < len(m.session.Questions)-1 {
					m.session.CurrentIndex++
					m = m.restoreQuestion()
				} else {
					return NewSummaryModel(m.session, m.title), nil
				}
			}

		case "p", "left":
			if m.session.CurrentIndex > 0 {
				m.session.CurrentIndex--
				m = m.restoreQuestion()
			}

		case "h":
			if !m.submitted {
				m.showHint = !m.showHint
			}
		}
	}

	return m, nil
}

func (m ExamModel) View() tea.View {
	q := m.session.Current()

	header := fmt.Sprintf("%s\n\nQuestion %d of %d\n\n%s\n",
		titleStyle.Render(m.title),
		m.session.CurrentIndex+1,
		len(m.session.Questions),
		q.Text,
	)

	// build a set of correct answer IDs for quick lookup
	correctIDs := make(map[string]bool)
	for _, id := range q.CorrectAnswers {
		correctIDs[id] = true
	}

	options := ""
	for i, opt := range q.Options {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if m.pending[opt.ID] {
			checked = "x"
		}

		row := fmt.Sprintf("%s [%s] %s", cursor, checked, opt.Text)

		if m.submitted {
			if correctIDs[opt.ID] {
				row = correctStyle.Render(row)
			} else if m.pending[opt.ID] {
				row = incorrectStyle.Render(row)
			}
		}

		options += row + "\n"
	}

	middle := ""
	if m.showHint && !m.submitted && q.Hint != "" {
		middle += fmt.Sprintf("\nHint: %s\n", q.Hint)
	}
	if m.submitted {
		middle += fmt.Sprintf("\n%s\n", q.Explanation)
	}

	footer := "space: select · enter: submit · h: hint · p: back · q: quit"

	body := header + "\n" + options + middle

	if m.width > 0 {
		v := tea.NewView(bodyStyle(m.width, m.height).Render(body) + "\n" + footerStyle(m.width).Render(footer))
		v.AltScreen = true
		return v
	}

	return tea.NewView(body + "\n" + footer + "\n")
}

// restoreQuestion resets cursor and showHint, then repopulates pending and
// submitted from the session so previously answered questions are shown correctly.
func (m ExamModel) restoreQuestion() ExamModel {
	q := m.session.Current()
	m.cursor = 0
	m.showHint = false
	m.submitted = m.session.IsAnswered(q.ID)
	m.pending = make(map[string]bool)
	for _, optID := range m.session.Answers[q.ID] {
		m.pending[optID] = true
	}
	return m
}
