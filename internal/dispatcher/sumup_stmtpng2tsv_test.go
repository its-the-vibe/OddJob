package dispatcher

import (
	"testing"

	"github.com/its-the-vibe/OddJob/internal/config"
)

func TestSumupStmtpng2tsvToPoppit(t *testing.T) {
	transformer := NewSumupStmtpng2tsvTransformer()
	cfg := config.PoppitConfig{
		Repo:   "its-the-vibe/OddJob",
		Branch: "refs/heads/main",
		Dir:    "/workspace",
		Type:   "odd:job",
	}

	msg, err := transformer.ToPoppit(TaskMessage{
		TaskName:  sumupStmtpng2tsvTaskName,
		InputFile: "/workspace/incoming/SumUp-Statement-Aug-26-?.png",
		Metadata: map[string]string{
			"stmtdate": "2026-08",
		},
	}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg.Dir != "${stmtpng2tsvDir}" {
		t.Fatalf("unexpected dir: %q", msg.Dir)
	}
	if len(msg.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(msg.Commands))
	}
	expectedCmd := `./stmtpng2tsv -output /workspace/incoming/SumUp-Statement-Aug-26.tsv /workspace/incoming/SumUp-Statement-Aug-26-?.png`
	if msg.Commands[0] != expectedCmd {
		t.Fatalf("unexpected command:\nexpected: %q\ngot:      %q", expectedCmd, msg.Commands[0])
	}
	if msg.Metadata["taskName"] != sumupStmtpng2tsvTaskName {
		t.Fatalf("expected metadata taskName %q, got %q", sumupStmtpng2tsvTaskName, msg.Metadata["taskName"])
	}
	if msg.Metadata["tsvFile"] != "/workspace/incoming/SumUp-Statement-Aug-26.tsv" {
		t.Fatalf("expected metadata tsvFile %q, got %q", "/workspace/incoming/SumUp-Statement-Aug-26.tsv", msg.Metadata["tsvFile"])
	}
	if msg.Metadata["stmtdate"] != "2026-08" {
		t.Fatalf("expected metadata stmtdate %q, got %q", "2026-08", msg.Metadata["stmtdate"])
	}
}

func TestSumupStmtpng2tsvToPoppitReturnsErrorWhenInputFileEmpty(t *testing.T) {
	transformer := NewSumupStmtpng2tsvTransformer()
	_, err := transformer.ToPoppit(TaskMessage{
		TaskName: sumupStmtpng2tsvTaskName,
	}, config.PoppitConfig{})
	if err == nil {
		t.Fatalf("expected error when inputFile is empty")
	}
}

func TestSumupStmtpng2tsvToPoppitSupportsInputFileWithSpaces(t *testing.T) {
	transformer := NewSumupStmtpng2tsvTransformer()

	msg, err := transformer.ToPoppit(TaskMessage{
		TaskName:  sumupStmtpng2tsvTaskName,
		InputFile: "/workspace/incoming/My SumUp-Statement-Aug-26-?.png",
	}, config.PoppitConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Dir != "${stmtpng2tsvDir}" {
		t.Fatalf("unexpected dir: %q", msg.Dir)
	}
	expectedCmd := `./stmtpng2tsv -output /workspace/incoming/My SumUp-Statement-Aug-26.tsv /workspace/incoming/My SumUp-Statement-Aug-26-?.png`
	if msg.Commands[0] != expectedCmd {
		t.Fatalf("unexpected command:\nexpected: %q\ngot:      %q", expectedCmd, msg.Commands[0])
	}
}

func TestSumupStmtpng2tsvFromPoppitChainsToStmt2redis(t *testing.T) {
	transformer := NewSumupStmtpng2tsvTransformer()

	task, ok, err := transformer.FromPoppit(PoppitOutput{
		StatusCode: 0,
		Metadata: map[string]string{
			"taskName": sumupStmtpng2tsvTaskName,
			"tsvFile":  "/workspace/incoming/SumUp-Statement-Aug-26.tsv",
			"stmtdate": "2026-08",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected chained task")
	}
	if task.TaskName != sumupStmt2redisTaskName {
		t.Fatalf("expected taskName %q, got %q", sumupStmt2redisTaskName, task.TaskName)
	}
	if task.InputFile != "/workspace/incoming/SumUp-Statement-Aug-26.tsv" {
		t.Fatalf("unexpected inputFile: %q", task.InputFile)
	}
	if task.Metadata["stmtdate"] != "2026-08" {
		t.Fatalf("expected metadata stmtdate %q, got %q", "2026-08", task.Metadata["stmtdate"])
	}
}

func TestSumupStmtpng2tsvFromPoppitReturnsErrorWhenTsvFileEmpty(t *testing.T) {
	transformer := NewSumupStmtpng2tsvTransformer()

	_, ok, err := transformer.FromPoppit(PoppitOutput{
		StatusCode: 0,
		Metadata: map[string]string{
			"taskName": sumupStmtpng2tsvTaskName,
		},
	})
	if err == nil {
		t.Fatalf("expected error when tsvFile is empty")
	}
	if ok {
		t.Fatalf("expected no chained task when tsvFile is empty")
	}
}

func TestSumupStmtpng2tsvFromPoppitReturnsNilWhenTaskNameMismatch(t *testing.T) {
	transformer := NewSumupStmtpng2tsvTransformer()

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

func TestSumupStmtpng2tsvFromPoppitReturnsNilWhenStatusCodeNonZero(t *testing.T) {
	transformer := NewSumupStmtpng2tsvTransformer()

	task, ok, err := transformer.FromPoppit(PoppitOutput{
		StatusCode: 1,
		Metadata: map[string]string{
			"taskName": sumupStmtpng2tsvTaskName,
			"tsvFile":  "/workspace/incoming/SumUp-Statement-Aug-26.tsv",
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
