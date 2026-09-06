package dispatcher

import (
	"fmt"
	"strings"

	"github.com/its-the-vibe/OddJob/internal/config"
)

const (
	sumupStmt2redisTaskName = "sumup:stmt2redis"
	sumupStmt2redisDir      = "${basedir}/${orgname}/stmt2redis"
)

type SumupStmt2redisTransformer struct{}

func NewSumupStmt2redisTransformer() *SumupStmt2redisTransformer {
	return &SumupStmt2redisTransformer{}
}

func (s *SumupStmt2redisTransformer) TaskName() string {
	return sumupStmt2redisTaskName
}

func (s *SumupStmt2redisTransformer) ToPoppit(task TaskMessage, cfg config.PoppitConfig) (PoppitMessage, error) {
	if task.TaskName != sumupStmt2redisTaskName {
		return PoppitMessage{}, fmt.Errorf("unsupported task for sumup stmt2redis transformer: %s", task.TaskName)
	}

	inputFile := strings.TrimSpace(task.InputFile)
	if inputFile == "" {
		return PoppitMessage{}, fmt.Errorf("inputFile is required for task %q", sumupStmt2redisTaskName)
	}

	metadata := make(map[string]string, len(task.Metadata)+1)
	for key, value := range task.Metadata {
		metadata[key] = value
	}
	metadata["taskName"] = sumupStmt2redisTaskName

	return PoppitMessage{
		Repo:   cfg.Repo,
		Branch: cfg.Branch,
		Type:   cfg.Type,
		Dir:    sumupStmt2redisDir,
		Commands: []string{
			fmt.Sprintf(`${stmt2redis} push -f %q -t sumup`, inputFile),
		},
		Metadata: metadata,
	}, nil
}

func (s *SumupStmt2redisTransformer) FromPoppit(output PoppitOutput) (*TaskMessage, bool, error) {
	return nil, false, nil
}
