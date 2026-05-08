package exam

type ExamFile struct {
	Meta      Meta       `json:"meta"`
	Questions []Question `json:"questions"`
}

type Meta struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
}

type Question struct {
	ID             string   `json:"id"`
	Type           string   `json:"type"` // "single" | "multi"
	Text           string   `json:"text"`
	Options        []Option `json:"options"`
	CorrectAnswers []string `json:"correct_answers"`
	Explanation    string   `json:"explanation"`
	Hints          []string `json:"hints,omitempty"`
}

type Option struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
