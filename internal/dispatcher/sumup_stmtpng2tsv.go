package dispatcher

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/its-the-vibe/OddJob/internal/config"
)

const sumupStmtpng2tsvTaskName = "sumup:stmtpng2tsv"
const sumupStmt2redisTaskName = "sumup:stmt2redis"

type SumupStmtpng2tsvTransformer struct{}

func NewSumupStmtpng2tsvTransformer() *SumupStmtpng2tsvTransformer {
	return &SumupStmtpng2tsvTransformer{}
}

func (s *SumupStmtpng2tsvTransformer) TaskName() string {
	return sumupStmtpng2tsvTaskName
}

func (s *SumupStmtpng2tsvTransformer) ToPoppit(task TaskMessage, cfg config.PoppitConfig) (PoppitMessage, error) {
	if task.TaskName != sumupStmtpng2tsvTaskName {
		return PoppitMessage{}, fmt.Errorf("unsupported task for sumup stmtpng2tsv transformer: %s", task.TaskName)
	}

	inputFile := strings.TrimSpace(task.InputFile)
	if inputFile == "" {
		return PoppitMessage{}, fmt.Errorf("inputFile is required for task %q", sumupStmtpng2tsvTaskName)
	}

	dir := "${stmtpng2tsvDir}"

	// Remove the extension and page number suffix (e.g., "/path/SumUp-Statement-Aug-26-1.png" -> "/path/SumUp-Statement-Aug-26")
	baseWithoutExt := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	baseWithoutPage := baseWithoutExt
	// Remove the page number suffix if it exists (e.g., "-1", "-2", etc.)
	if len(baseWithoutExt) > 2 && baseWithoutExt[len(baseWithoutExt)-2] == '-' {
		if baseWithoutExt[len(baseWithoutExt)-1] >= '0' && baseWithoutExt[len(baseWithoutExt)-1] <= '9' {
			baseWithoutPage = baseWithoutExt[:len(baseWithoutExt)-2]
		}
	}

	outputFile := baseWithoutPage + ".tsv"
	inputPattern := baseWithoutPage + "-?.png"

	metadata := make(map[string]string, len(task.Metadata)+2)
	for key, value := range task.Metadata {
		metadata[key] = value
	}
	metadata["taskName"] = sumupStmtpng2tsvTaskName
	metadata["tsvFile"] = outputFile

	return PoppitMessage{
		Repo:   cfg.Repo,
		Branch: cfg.Branch,
		Type:   cfg.Type,
		Dir:    dir,
		Commands: []string{
			fmt.Sprintf(`echo ./stmtpng2tsv -output %q %q`, outputFile, inputPattern),
		},
		Metadata: metadata,
	}, nil
}

func (s *SumupStmtpng2tsvTransformer) FromPoppit(output PoppitOutput) (*TaskMessage, bool, error) {
	if output.Metadata["taskName"] != sumupStmtpng2tsvTaskName {
		return nil, false, nil
	}
	if output.StatusCode != 0 {
		return nil, false, nil
	}

	tsvFile := strings.TrimSpace(output.Metadata["tsvFile"])
	if tsvFile == "" {
		return nil, false, fmt.Errorf("missing tsvFile metadata for task %q", sumupStmtpng2tsvTaskName)
	}

	metadata := map[string]string{}
	if stmtdate := strings.TrimSpace(output.Metadata["stmtdate"]); stmtdate != "" {
		metadata["stmtdate"] = stmtdate
	}

	return &TaskMessage{
		TaskName:  sumupStmt2redisTaskName,
		InputFile: tsvFile,
		Metadata:  metadata,
	}, true, nil
}
