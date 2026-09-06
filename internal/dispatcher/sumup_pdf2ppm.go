package dispatcher

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/its-the-vibe/OddJob/internal/config"
)

const sumupPdf2ppmTaskName = "sumup:pdf2ppm"

type SumupPdf2ppmTransformer struct{}

func NewSumupPdf2ppmTransformer() *SumupPdf2ppmTransformer {
	return &SumupPdf2ppmTransformer{}
}

func (s *SumupPdf2ppmTransformer) TaskName() string {
	return sumupPdf2ppmTaskName
}

func (s *SumupPdf2ppmTransformer) ToPoppit(task TaskMessage, cfg config.PoppitConfig) (PoppitMessage, error) {
	if task.TaskName != sumupPdf2ppmTaskName {
		return PoppitMessage{}, fmt.Errorf("unsupported task for sumup pdftoppm transformer: %s", task.TaskName)
	}

	inputFile := strings.TrimSpace(task.InputFile)
	if inputFile == "" {
		return PoppitMessage{}, fmt.Errorf("inputFile is required for task %q", sumupPdf2ppmTaskName)
	}

	dir := filepath.Dir(inputFile)
	baseFile := filepath.Base(inputFile)
	outputPrefix := strings.TrimSuffix(baseFile, filepath.Ext(baseFile))
	if outputPrefix == "" {
		return PoppitMessage{}, fmt.Errorf("invalid inputFile for task %q: %q", sumupPdf2ppmTaskName, inputFile)
	}
	pngFile := filepath.Join(dir, fmt.Sprintf("%s-?.png", outputPrefix))

	metadata := make(map[string]string, len(task.Metadata)+2)
	for key, value := range task.Metadata {
		metadata[key] = value
	}
	metadata["taskName"] = sumupPdf2ppmTaskName
	metadata["pngFile"] = pngFile

	return PoppitMessage{
		Repo:   cfg.Repo,
		Branch: cfg.Branch,
		Type:   cfg.Type,
		Dir:    dir,
		Commands: []string{
			fmt.Sprintf(`pdftoppm -png -r 300 %q %q`, baseFile, outputPrefix),
		},
		Metadata: metadata,
	}, nil
}

func (s *SumupPdf2ppmTransformer) FromPoppit(output PoppitOutput) (*TaskMessage, bool, error) {
	if output.Metadata["taskName"] != sumupPdf2ppmTaskName {
		return nil, false, nil
	}
	if output.StatusCode != 0 {
		return nil, false, nil
	}

	pngFile := strings.TrimSpace(output.Metadata["pngFile"])
	if pngFile == "" {
		return nil, false, fmt.Errorf("missing pngFile metadata for task %q", sumupPdf2ppmTaskName)
	}

	metadata := map[string]string{}
	if stmtdate := strings.TrimSpace(output.Metadata["stmtdate"]); stmtdate != "" {
		metadata["stmtdate"] = stmtdate
	}

	return &TaskMessage{
		TaskName:  sumupStmtpng2tsvTaskName,
		InputFile: pngFile,
		Metadata:  metadata,
	}, true, nil
}
