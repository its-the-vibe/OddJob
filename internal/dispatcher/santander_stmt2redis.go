package dispatcher

import (
	"fmt"
	"strings"

	"github.com/its-the-vibe/OddJob/internal/config"
)

const (
	santanderStmt2redisTaskName = "santander:stmt2redis"
	santanderStmt2redisDir      = "${basedir}/${orgname}/stmt2redis"
)

type SantanderStmt2redisTransformer struct{}

func NewSantanderStmt2redisTransformer() *SantanderStmt2redisTransformer {
	return &SantanderStmt2redisTransformer{}
}

func (s *SantanderStmt2redisTransformer) TaskName() string {
	return santanderStmt2redisTaskName
}

func (s *SantanderStmt2redisTransformer) ToPoppit(task TaskMessage, cfg config.PoppitConfig) (PoppitMessage, error) {
	if task.TaskName != santanderStmt2redisTaskName {
		return PoppitMessage{}, fmt.Errorf("unsupported task for santander stmt2redis transformer: %s", task.TaskName)
	}

	inputFile := strings.TrimSpace(task.InputFile)
	if inputFile == "" {
		return PoppitMessage{}, fmt.Errorf("inputFile is required for task %q", santanderStmt2redisTaskName)
	}

	metadata := make(map[string]string, len(task.Metadata)+1)
	for key, value := range task.Metadata {
		metadata[key] = value
	}
	metadata["taskName"] = santanderStmt2redisTaskName

	return PoppitMessage{
		Repo:   cfg.Repo,
		Branch: cfg.Branch,
		Type:   cfg.Type,
		Dir:    santanderStmt2redisDir,
		Commands: []string{
			fmt.Sprintf(`${stmt2redis} push -f %q -t santander`, inputFile),
		},
		Metadata: metadata,
	}, nil
}

func (s *SantanderStmt2redisTransformer) FromPoppit(output PoppitOutput) (*TaskMessage, bool, error) {
	return nil, false, nil
}
