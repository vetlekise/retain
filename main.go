// Package main is the entry point for Retain.
package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/vetlekise/retain/internal/exam"
	"github.com/vetlekise/retain/internal/session"
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

	session.New(ef)
}
