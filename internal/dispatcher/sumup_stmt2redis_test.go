package dispatcher

import (
	"testing"

	"github.com/its-the-vibe/OddJob/internal/config"
)

func TestSumupStmt2redisToPoppit(t *testing.T) {
	transformer := NewSumupStmt2redisTransformer()
	cfg := config.PoppitConfig{
		Repo:   "its-the-vibe/OddJob",
		Branch: "refs/heads/main",
		Dir:    "/workspace",
		Type:   "odd:job",
	}

	msg, err := transformer.ToPoppit(TaskMessage{
		TaskName:  sumupStmt2redisTaskName,
		InputFile: "/workspace/incoming/SumUp-Statement-Aug-26.tsv",
		Metadata: map[string]string{
			"stmtdate": "2026-08",
		},
	}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg.Dir != sumupStmt2redisDir {
		t.Fatalf("expected dir %q, got %q", sumupStmt2redisDir, msg.Dir)
	}
	if len(msg.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(msg.Commands))
	}
	if msg.Commands[0] != `${stmt2redis} push -f "/workspace/incoming/SumUp-Statement-Aug-26.tsv" -t sumup` {
		t.Fatalf("unexpected command: %q", msg.Commands[0])
	}
	if msg.Metadata["taskName"] != sumupStmt2redisTaskName {
		t.Fatalf("expected metadata taskName %q, got %q", sumupStmt2redisTaskName, msg.Metadata["taskName"])
	}
	if msg.Metadata["stmtdate"] != "2026-08" {
		t.Fatalf("expected metadata stmtdate %q, got %q", "2026-08", msg.Metadata["stmtdate"])
	}
}

func TestSumupStmt2redisToPoppitReturnsErrorWhenInputFileEmpty(t *testing.T) {
	transformer := NewSumupStmt2redisTransformer()
	_, err := transformer.ToPoppit(TaskMessage{
		TaskName: sumupStmt2redisTaskName,
	}, config.PoppitConfig{})
	if err == nil {
		t.Fatalf("expected error when inputFile is empty")
	}
}

func TestSumupStmt2redisToPoppitPreservesMetadata(t *testing.T) {
	transformer := NewSumupStmt2redisTransformer()

	msg, err := transformer.ToPoppit(TaskMessage{
		TaskName:  sumupStmt2redisTaskName,
		InputFile: "/workspace/incoming/SumUp-Statement-Aug-26.tsv",
		Metadata: map[string]string{
			"stmtdate": "2026-08",
			"custom":   "value",
		},
	}, config.PoppitConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg.Metadata["stmtdate"] != "2026-08" {
		t.Fatalf("expected metadata stmtdate %q, got %q", "2026-08", msg.Metadata["stmtdate"])
	}
	if msg.Metadata["custom"] != "value" {
		t.Fatalf("expected metadata custom %q, got %q", "value", msg.Metadata["custom"])
	}
}

func TestSumupStmt2redisFromPoppit(t *testing.T) {
	transformer := NewSumupStmt2redisTransformer()

	task, ok, err := transformer.FromPoppit(PoppitOutput{
		StatusCode: 0,
		Metadata: map[string]string{
			"taskName": sumupStmt2redisTaskName,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("expected no chained task")
	}
	if task != nil {
		t.Fatalf("expected nil task")
	}
}
