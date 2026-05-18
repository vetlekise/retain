package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/vetlekise/retain/internal/session"
)

type ExamModel struct {
	session   *session.Session
	submitted bool            // has the current question been submitted?
	pending   map[string]bool // option ID → selected
	cursor    int
	showHint  bool
}

func NewExamModel(s *session.Session) ExamModel {
	return ExamModel{
		session: s,
		pending: make(map[string]bool),
	}
}

func (m ExamModel) Init() tea.Cmd {
	return nil
}

func (m ExamModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
					m.pending = make(map[string]bool)
					m.submitted = false
					m.cursor = 0
					m.showHint = false
				}
			}

		case "p", "left":
			if m.session.CurrentIndex > 0 {
				m.session.CurrentIndex--
				m.pending = make(map[string]bool)
				m.submitted = false
				m.cursor = 0
				m.showHint = false
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

	s := fmt.Sprintf("Question %d of %d\n\n%s\n\n",
		m.session.CurrentIndex+1,
		len(m.session.Questions),
		q.Text,
	)

	for i, opt := range q.Options {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if m.pending[opt.ID] {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, opt.Text)
	}

	if m.showHint && !m.submitted && q.Hint != "" {
		s += fmt.Sprintf("\nHint: %s\n", q.Hint)
	}

	if m.submitted {
		s += fmt.Sprintf("\n%s\n", q.Explanation)
	}

	s += "\nspace: select · enter: submit · h: hint · p: back · q: quit\n"

	return tea.NewView(s)
}
