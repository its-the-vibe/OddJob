package dispatcher

import "fmt"

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
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Type       string                 `json:"type"`
	Command    string                 `json:"command"`
	Output     string                 `json:"output"`
	Stderr     string                 `json:"stderr"`
	StatusCode int                    `json:"status_code"`
}

// metadataString returns the metadata value for key as a string, converting
// non-string JSON values (numbers, booleans) to their string form.
func metadataString(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}
