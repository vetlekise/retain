package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/vetlekise/retain/internal/session"
)

type ExamModel struct {
	session   *session.Session
	title     string
	submitted bool            // has the current question been submitted?
	pending   map[string]bool // option ID → selected
	cursor    int
	showHint  bool
	width     int
	height    int
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

	titleStyle := lipgloss.NewStyle().Bold(true)
	header := fmt.Sprintf("%s\n\nQuestion %d of %d\n\n%s\n",
		titleStyle.Render(m.title),
		m.session.CurrentIndex+1,
		len(m.session.Questions),
		q.Text,
	)

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

		options += fmt.Sprintf("%s [%s] %s\n", cursor, checked, opt.Text)
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
		style := lipgloss.NewStyle().Width(m.width).Height(m.height-1).Padding(1, 2)
		footerStyle := lipgloss.NewStyle().Width(m.width)
		v := tea.NewView(style.Render(body) + "\n" + footerStyle.Render(footer))
		v.AltScreen = true
		return v
	}

	return tea.NewView(body + "\n" + footer + "\n")
}
