// Package main is the entry point for Retain.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/vetlekise/retain/internal/exam"
	"github.com/vetlekise/retain/internal/session"
	"github.com/vetlekise/retain/internal/ui"
)

var examPath string

func init() {
	flag.StringVar(&examPath, "exam", "exam.json", "Path to your .json defined exam file.")
	flag.Parse()
}

func main() {

	data, err := exam.ReadExam(examPath)
	if err != nil {
		slog.Error("failed to read exam", "err", err)
		os.Exit(1)
	}

	ef, err := exam.LoadExam(data)
	if err != nil {
		slog.Error("failed to load exam", "err", err)
		os.Exit(1)
	}

	sess := session.New(ef)

	p := tea.NewProgram(ui.NewExamModel(sess, ef.Meta.Title))
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
