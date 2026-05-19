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

// New returns a Session with Questions from ef, CurrentIndex 0,
// an initialized empty Answers map, and StartTime set to now
func New(ef exam.ExamFile) *Session {
	return &Session{
		Questions:    ef.Questions,
		CurrentIndex: 0,
		Answers:      make(map[string][]string),
		StartTime:    time.Now(),
	}
}

// Current returns the current question
func (s *Session) Current() exam.Question {
	return s.Questions[s.CurrentIndex]
}

// Answer records/overwrites an answer
func (s *Session) Answer(questionID string, optionIDs []string) {
	s.Answers[questionID] = optionIDs
}

// IsAnswered looks up questionID in s.Answers
//
// Returns true if it exists and has at least one option
func (s *Session) IsAnswered(questionID string) bool {
	optionIDs, ok := s.Answers[questionID]
	return ok && len(optionIDs) > 0
}

// equalUnordered reports whether a and b contain the same elements regardless of order.
func equalUnordered(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	// frequency map: increment for each element in a, decrement for each in b.
	// If any count goes negative, b contains an element not in a.
	counts := make(map[string]int)
	for _, v := range a {
		counts[v]++
	}
	for _, v := range b {
		counts[v]--
		if counts[v] < 0 {
			return false
		}
	}
	return true
}

// IsCorrect reports whether the recorded answer for questionID matches the correct answers.
func (s *Session) IsCorrect(questionID string) bool {
	answer, ok := s.Answers[questionID]
	if !ok || len(answer) == 0 {
		return false
	}

	for _, question := range s.Questions {
		if question.ID == questionID {
			return equalUnordered(answer, question.CorrectAnswers)
		}
	}

	return false
}

// Score returns the number of correctly answered questions and the total number of questions.
func (s *Session) Score() (correct, total int) {
	total = len(s.Questions)

	for _, question := range s.Questions {
		if s.IsCorrect(question.ID) {
			correct++
		}
	}

	return correct, total
}

// Elapsed returns the time elapsed since the session started.
func (s *Session) Elapsed() time.Duration {
	return time.Since(s.StartTime)
}

// Reset clears all answers and restarts the timer, keeping the same questions.
func (s *Session) Reset() {
	s.CurrentIndex = 0
	s.Answers = make(map[string][]string)
	s.StartTime = time.Now()
}
