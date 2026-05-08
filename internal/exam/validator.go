package exam

import "fmt"

// Validate parses an unmarshaled JSON (ExamFile) and returns a list of errors (errs),
// so users see all problems in one run.
func Validate(ef ExamFile) []error {
	var errs []error
	seenIDs := make(map[string]bool)

	for i, q := range ef.Questions {
		if q.ID == "" {
			errs = append(errs, fmt.Errorf("question %d: missing id", i))
		}
		if q.Text == "" {
			errs = append(errs, fmt.Errorf("question %s: missing text", q.ID))
		}
		if len(q.Options) == 0 {
			errs = append(errs, fmt.Errorf("question %s: no options", q.ID))
		}
		if len(q.CorrectAnswers) == 0 {
			errs = append(errs, fmt.Errorf("question %s: no correct answers", q.ID))
		}

		// duplicate ID check
		if seenIDs[q.ID] {
			errs = append(errs, fmt.Errorf("question %s: duplicate id", q.ID))
		}
		seenIDs[q.ID] = true

		// correct answers reference valid option IDs
		optionIDs := make(map[string]bool)
		for _, o := range q.Options {
			optionIDs[o.ID] = true
		}
		for _, ca := range q.CorrectAnswers {
			if !optionIDs[ca] {
				errs = append(errs, fmt.Errorf("question %s: correct answer %q does not match any option", q.ID, ca))
			}
		}
	}

	return errs
}
