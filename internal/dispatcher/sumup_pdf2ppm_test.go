package dispatcher

import (
	"testing"

	"github.com/its-the-vibe/OddJob/internal/config"
)

func TestSumupPdftoppmToPoppit(t *testing.T) {
	transformer := NewSumupPdf2ppmTransformer()
	cfg := config.PoppitConfig{
		Repo:   "its-the-vibe/OddJob",
		Branch: "refs/heads/main",
		Dir:    "/workspace",
		Type:   "odd:job",
	}

	msg, err := transformer.ToPoppit(TaskMessage{
		TaskName:  sumupPdf2ppmTaskName,
		InputFile: "/workspace/incoming/SumUp-Statement-Aug-26.pdf",
		Metadata:  map[string]string{},
	}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg.Dir != "/workspace/incoming" {
		t.Fatalf("unexpected dir: %q", msg.Dir)
	}
	if len(msg.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(msg.Commands))
	}
	if msg.Commands[0] != `pdftoppm -png -r 300 "SumUp-Statement-Aug-26.pdf" "SumUp-Statement-Aug-26"` {
		t.Fatalf("unexpected command: %q", msg.Commands[0])
	}
	if msg.Metadata["taskName"] != sumupPdf2ppmTaskName {
		t.Fatalf("expected metadata taskName %q, got %q", sumupPdf2ppmTaskName, msg.Metadata["taskName"])
	}
	if msg.Metadata["pngFile"] != "/workspace/incoming/SumUp-Statement-Aug-26-?.png" {
		t.Fatalf("expected metadata pngFile %q, got %q", "/workspace/incoming/SumUp-Statement-Aug-26-?.png", msg.Metadata["pngFile"])
	}
}

func TestSumupPdftoppmToPoppitReturnsErrorWhenInputFileEmpty(t *testing.T) {
	transformer := NewSumupPdf2ppmTransformer()
	_, err := transformer.ToPoppit(TaskMessage{
		TaskName: sumupPdf2ppmTaskName,
	}, config.PoppitConfig{})
	if err == nil {
		t.Fatalf("expected error when inputFile is empty")
	}
}

func TestSumupPdftoppmToPoppitSupportsInputFileWithSpaces(t *testing.T) {
	transformer := NewSumupPdf2ppmTransformer()

	msg, err := transformer.ToPoppit(TaskMessage{
		TaskName:  sumupPdf2ppmTaskName,
		InputFile: "/workspace/incoming/My SumUp-Statement-Aug-26.pdf",
	}, config.PoppitConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Dir != "/workspace/incoming" {
		t.Fatalf("unexpected dir: %q", msg.Dir)
	}
	if msg.Commands[0] != `pdftoppm -png -r 300 "My SumUp-Statement-Aug-26.pdf" "My SumUp-Statement-Aug-26"` {
		t.Fatalf("unexpected command: %q", msg.Commands[0])
	}
}

func TestSumupPdftoppmFromPoppitChainsToStmtpng2tsv(t *testing.T) {
	transformer := NewSumupPdf2ppmTransformer()

	task, ok, err := transformer.FromPoppit(PoppitOutput{
		StatusCode: 0,
		Metadata: map[string]string{
			"taskName": sumupPdf2ppmTaskName,
			"pngFile":  "/workspace/incoming/SumUp-Statement-Aug-26-1.png",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected chained task")
	}
	if task.TaskName != sumupStmtpng2tsvTaskName {
		t.Fatalf("expected taskName %q, got %q", sumupStmtpng2tsvTaskName, task.TaskName)
	}
	if task.InputFile != "/workspace/incoming/SumUp-Statement-Aug-26-1.png" {
		t.Fatalf("unexpected inputFile: %q", task.InputFile)
	}
}

func TestSumupPdftoppmFromPoppitIgnoresMetadataFromPrevious(t *testing.T) {
	transformer := NewSumupPdf2ppmTransformer()

	task, ok, err := transformer.FromPoppit(PoppitOutput{
		StatusCode: 0,
		Metadata: map[string]string{
			"taskName": sumupPdf2ppmTaskName,
			"pngFile":  "/workspace/incoming/SumUp-Statement-Aug-26-1.png",
			"stmtdate": "2026-08",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected chained task")
	}
	if task.Metadata["stmtdate"] != "2026-08" {
		t.Fatalf("expected stmtdate to be preserved in metadata")
	}
}

func TestSumupPdftoppmFromPoppitReturnsNilWhenTaskNameMismatch(t *testing.T) {
	transformer := NewSumupPdf2ppmTransformer()

	task, ok, err := transformer.FromPoppit(PoppitOutput{
		StatusCode: 0,
		Metadata: map[string]string{
			"taskName": "different:task",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("expected no chained task for mismatched taskName")
	}
	if task != nil {
		t.Fatalf("expected nil task for mismatched taskName")
	}
}

func TestSumupPdftoppmFromPoppitReturnsNilWhenStatusCodeNonZero(t *testing.T) {
	transformer := NewSumupPdf2ppmTransformer()

	task, ok, err := transformer.FromPoppit(PoppitOutput{
		StatusCode: 1,
		Metadata: map[string]string{
			"taskName": sumupPdf2ppmTaskName,
			"pngFile":  "/workspace/incoming/SumUp-Statement-Aug-26-1.png",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("expected no chained task when status code is non-zero")
	}
	if task != nil {
		t.Fatalf("expected nil task when status code is non-zero")
	}
}
