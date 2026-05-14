package dispatcher

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/its-the-vibe/OddJob/internal/config"
)

const santanderStmtpng2tsvTaskName = "santander:stmtpng2tsv"

type SantanderStmtpng2tsvTransformer struct{}

func NewSantanderStmtpng2tsvTransformer() *SantanderStmtpng2tsvTransformer {
	return &SantanderStmtpng2tsvTransformer{}
}

func (s *SantanderStmtpng2tsvTransformer) TaskName() string {
	return santanderStmtpng2tsvTaskName
}

func (s *SantanderStmtpng2tsvTransformer) ToPoppit(task TaskMessage, cfg config.PoppitConfig) (PoppitMessage, error) {
	if task.TaskName != santanderStmtpng2tsvTaskName {
		return PoppitMessage{}, fmt.Errorf("unsupported task for santander stmtpng2tsv transformer: %s", task.TaskName)
	}

	inputFile := strings.TrimSpace(task.InputFile)
	if inputFile == "" {
		return PoppitMessage{}, fmt.Errorf("inputFile is required for task %q", santanderStmtpng2tsvTaskName)
	}

	dir := filepath.Dir(inputFile)
	inputBase := filepath.Base(inputFile)
	outputBase := strings.TrimSuffix(inputBase, filepath.Ext(inputBase))
	if outputBase == "" {
		return PoppitMessage{}, fmt.Errorf("invalid inputFile for task %q: %q", santanderStmtpng2tsvTaskName, inputFile)
	}
	outputFile := filepath.Join(dir, outputBase+".tsv")

	metadata := make(map[string]string, len(task.Metadata)+2)
	for key, value := range task.Metadata {
		metadata[key] = value
	}
	metadata["taskName"] = santanderStmtpng2tsvTaskName
	metadata["tsvFile"] = outputFile

	return PoppitMessage{
		Repo:   cfg.Repo,
		Branch: cfg.Branch,
		Type:   cfg.Type,
		Dir:    dir,
		Commands: []string{
			fmt.Sprintf(`${stmtpng2tsv} -input %q -output %q`, inputBase, outputBase+".tsv"),
		},
		Metadata: metadata,
	}, nil
}

func (s *SantanderStmtpng2tsvTransformer) FromPoppit(output PoppitOutput) (*TaskMessage, bool, error) {
	if output.Metadata["taskName"] != santanderStmtpng2tsvTaskName {
		return nil, false, nil
	}
	if output.StatusCode != 0 {
		return nil, false, nil
	}

	tsvFile := strings.TrimSpace(output.Metadata["tsvFile"])
	if tsvFile == "" {
		return nil, false, fmt.Errorf("missing tsvFile metadata for task %q", santanderStmtpng2tsvTaskName)
	}

	metadata := map[string]string{}
	if stmtdate := strings.TrimSpace(output.Metadata["stmtdate"]); stmtdate != "" {
		metadata["stmtdate"] = stmtdate
	}

	return &TaskMessage{
		TaskName:  santanderStmt2redisTaskName,
		InputFile: tsvFile,
		Metadata:  metadata,
	}, true, nil
}
