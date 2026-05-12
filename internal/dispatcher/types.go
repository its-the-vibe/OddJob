package dispatcher

type TaskMessage struct {
	TaskName  string            `json:"taskName"`
	InputFile string            `json:"inputFile,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type PoppitMessage struct {
	Repo     string            `json:"repo"`
	Branch   string            `json:"branch"`
	Type     string            `json:"type"`
	Dir      string            `json:"dir"`
	Commands []string          `json:"commands"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type PoppitOutput struct {
	Metadata   map[string]string `json:"metadata,omitempty"`
	Type       string            `json:"type"`
	Command    string            `json:"command"`
	Output     string            `json:"output"`
	Stderr     string            `json:"stderr"`
	StatusCode int               `json:"status_code"`
}
