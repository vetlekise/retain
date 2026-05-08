package session

import (
	"time"

	"github.com/vetlekise/retain/internal/exam"
)

type Session struct {
	Questions    []exam.Question
	CurrentIndex int
	Answers      map[string][]string // question ID > selected option IDs
	StartTime    time.Time
}

func New(ef exam.ExamFile) *Session

func Current() exam.Question

func Answer(questionID string, optionIDs []string)

func IsAnswered(questionID string) bool

func IsCorrect(questionID string) bool

func Score() (correct, total int)

func Elapsed() time.Duration
