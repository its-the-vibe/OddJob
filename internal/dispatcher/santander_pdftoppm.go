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
	pngFile := filepath.Join(dir, fmt.Sprintf("%s-2.png", outputPrefix))

	metadata := make(map[string]string, len(task.Metadata)+2)
	for key, value := range task.Metadata {
		metadata[key] = value
	}
	metadata["taskName"] = santanderPdftoppmTaskName
	metadata["pngFile"] = pngFile

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
	if output.Metadata["taskName"] != santanderPdftoppmTaskName {
		return nil, false, nil
	}
	if output.StatusCode != 0 {
		return nil, false, nil
	}

	pngFile := strings.TrimSpace(output.Metadata["pngFile"])
	if pngFile == "" {
		return nil, false, fmt.Errorf("missing pngFile metadata for task %q", santanderPdftoppmTaskName)
	}

	metadata := map[string]string{}
	if stmtdate := strings.TrimSpace(output.Metadata["stmtdate"]); stmtdate != "" {
		metadata["stmtdate"] = stmtdate
	}

	return &TaskMessage{
		TaskName:  santanderStmtpng2tsvTaskName,
		InputFile: pngFile,
		Metadata:  metadata,
	}, true, nil
}
