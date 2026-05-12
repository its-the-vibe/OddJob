package dispatcher

import (
	"fmt"
	"strings"

	"github.com/its-the-vibe/OddJob/internal/config"
)

const santanderStmtdateTaskName = "santander:stmtdate"

type SantanderStmtdateTransformer struct{}

func NewSantanderStmtdateTransformer() *SantanderStmtdateTransformer {
	return &SantanderStmtdateTransformer{}
}

func (s *SantanderStmtdateTransformer) TaskName() string {
	return santanderStmtdateTaskName
}

func (s *SantanderStmtdateTransformer) ToPoppit(task TaskMessage, cfg config.PoppitConfig) (PoppitMessage, error) {
	if task.TaskName != santanderStmtdateTaskName {
		return PoppitMessage{}, fmt.Errorf("unsupported task for santander stmtdate transformer: %s", task.TaskName)
	}

	filename := strings.TrimSpace(task.InputFile)
	if filename == "" {
		return PoppitMessage{}, fmt.Errorf("inputFile is required for task %q", santanderStmtdateTaskName)
	}

	return PoppitMessage{
		Repo:     cfg.Repo,
		Branch:   cfg.Branch,
		Type:     cfg.Type,
		Dir:      cfg.Dir,
		Commands: []string{fmt.Sprintf("stmtdate -rename %q", filename)},
		Metadata: map[string]string{
			"taskName": santanderStmtdateTaskName,
		},
	}, nil
}

func (s *SantanderStmtdateTransformer) FromPoppit(output PoppitOutput) (*TaskMessage, bool, error) {
	return nil, false, nil
}
