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
		Commands: []string{fmt.Sprintf("${stmtdate} -rename %q", filename)},
		Metadata: map[string]string{
			"taskName": santanderStmtdateTaskName,
		},
	}, nil
}

func (s *SantanderStmtdateTransformer) FromPoppit(output PoppitOutput) (*TaskMessage, bool, error) {
	if output.Metadata["taskName"] != santanderStmtdateTaskName {
		return nil, false, nil
	}
	if output.StatusCode != 0 {
		return nil, false, nil
	}

	stmtdate, renamedFile, err := parseStmtdateOutput(output.Output)
	if err != nil {
		return nil, false, err
	}

	return &TaskMessage{
		TaskName:  santanderPdftoppmTaskName,
		InputFile: renamedFile,
		Metadata: map[string]string{
			"stmtdate": stmtdate,
		},
	}, true, nil
}

func parseStmtdateOutput(raw string) (string, string, error) {
	lines := strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n")
	nonEmpty := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		nonEmpty = append(nonEmpty, trimmed)
	}
	if len(nonEmpty) < 2 {
		return "", "", fmt.Errorf("unexpected stmtdate output: expected at least 2 non-empty lines")
	}

	stmtdate := nonEmpty[0]
	if stmtdate == "" {
		return "", "", fmt.Errorf("unexpected stmtdate output: missing statement date")
	}

	const renamedPrefix = "renamed:"
	renamedLine := nonEmpty[1]
	if !strings.HasPrefix(renamedLine, renamedPrefix) {
		return "", "", fmt.Errorf("unexpected stmtdate output: second line must start with %q", renamedPrefix)
	}
	renamedFile := strings.TrimSpace(strings.TrimPrefix(renamedLine, renamedPrefix))
	if renamedFile == "" {
		return "", "", fmt.Errorf("unexpected stmtdate output: missing renamed file path")
	}

	return stmtdate, renamedFile, nil
}
