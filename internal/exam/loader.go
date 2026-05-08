package exam

import (
	"encoding/json"
	"fmt"
)

// LoadExam loads
func LoadExam(file []byte) (ExamFile, error) {
	var ef ExamFile

	if err := json.Unmarshal(file, &ef); err != nil {
		return ExamFile{}, fmt.Errorf("failed to unmarshal: %w", err)
	}

	if err := Validate(ef); err != nil {
		return ExamFile{}, fmt.Errorf("validation failed: %v", err)
	}

	return ef, nil
}
