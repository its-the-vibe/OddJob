package dispatcher

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/its-the-vibe/OddJob/internal/config"
)

const santanderPdftoppmTaskName = "santander:pdftoppm"

type SantanderPdftoppmTransformer struct{}

func NewSantanderPdftoppmTransformer() *SantanderPdftoppmTransformer {
	return &SantanderPdftoppmTransformer{}
}

func (s *SantanderPdftoppmTransformer) TaskName() string {
	return santanderPdftoppmTaskName
}

func (s *SantanderPdftoppmTransformer) ToPoppit(task TaskMessage, cfg config.PoppitConfig) (PoppitMessage, error) {
	if task.TaskName != santanderPdftoppmTaskName {
		return PoppitMessage{}, fmt.Errorf("unsupported task for santander pdftoppm transformer: %s", task.TaskName)
	}

	inputFile := strings.TrimSpace(task.InputFile)
	if inputFile == "" {
		return PoppitMessage{}, fmt.Errorf("inputFile is required for task %q", santanderPdftoppmTaskName)
	}

	dir := filepath.Dir(inputFile)
	baseFile := filepath.Base(inputFile)
	outputPrefix := strings.TrimSuffix(baseFile, filepath.Ext(baseFile))
	if outputPrefix == "" {
		return PoppitMessage{}, fmt.Errorf("invalid inputFile for task %q: %q", santanderPdftoppmTaskName, inputFile)
	}

	metadata := make(map[string]string, len(task.Metadata)+1)
	for key, value := range task.Metadata {
		metadata[key] = value
	}
	metadata["taskName"] = santanderPdftoppmTaskName

	return PoppitMessage{
		Repo:   cfg.Repo,
		Branch: cfg.Branch,
		Type:   cfg.Type,
		Dir:    dir,
		Commands: []string{
			fmt.Sprintf(`pdftoppm -png -r 300 %q -f 2 %q`, baseFile, outputPrefix),
		},
		Metadata: metadata,
	}, nil
}

func (s *SantanderPdftoppmTransformer) FromPoppit(output PoppitOutput) (*TaskMessage, bool, error) {
	return nil, false, nil
}
