package dispatcher

import (
	"testing"

	"github.com/its-the-vibe/OddJob/internal/config"
)

func TestSantanderStmtpng2tsvToPoppit(t *testing.T) {
	transformer := NewSantanderStmtpng2tsvTransformer()
	cfg := config.PoppitConfig{
		Repo:   "its-the-vibe/OddJob",
		Branch: "refs/heads/main",
		Dir:    "/workspace",
		Type:   "odd:job",
	}

	msg, err := transformer.ToPoppit(TaskMessage{
		TaskName:  santanderStmtpng2tsvTaskName,
		InputFile: "/workspace/incoming/File-2026-03-2.png",
		Metadata: map[string]string{
			"stmtdate": "2026-03",
		},
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
	if msg.Commands[0] != `${stmtpng2tsv} -input "File-2026-03-2.png" -output "File-2026-03-2.tsv"` {
		t.Fatalf("unexpected command: %q", msg.Commands[0])
	}
	if msg.Metadata["taskName"] != santanderStmtpng2tsvTaskName {
		t.Fatalf("expected metadata taskName %q, got %q", santanderStmtpng2tsvTaskName, msg.Metadata["taskName"])
	}
	if msg.Metadata["tsvFile"] != "/workspace/incoming/File-2026-03-2.tsv" {
		t.Fatalf("expected metadata tsvFile %q, got %q", "/workspace/incoming/File-2026-03-2.tsv", msg.Metadata["tsvFile"])
	}
}

func TestSantanderStmtpng2tsvToPoppitReturnsErrorWhenInputFileEmpty(t *testing.T) {
	transformer := NewSantanderStmtpng2tsvTransformer()
	_, err := transformer.ToPoppit(TaskMessage{
		TaskName: santanderStmtpng2tsvTaskName,
	}, config.PoppitConfig{})
	if err == nil {
		t.Fatalf("expected error when inputFile is empty")
	}
}

func TestSantanderStmtpng2tsvFromPoppitChainsToStmt2redis(t *testing.T) {
	transformer := NewSantanderStmtpng2tsvTransformer()

	task, ok, err := transformer.FromPoppit(PoppitOutput{
		StatusCode: 0,
		Metadata: map[string]string{
			"taskName": santanderStmtpng2tsvTaskName,
			"tsvFile":  "/workspace/incoming/File-2026-03-2.tsv",
			"stmtdate": "2026-03",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected chained task")
	}
	if task.TaskName != santanderStmt2redisTaskName {
		t.Fatalf("expected taskName %q, got %q", santanderStmt2redisTaskName, task.TaskName)
	}
	if task.InputFile != "/workspace/incoming/File-2026-03-2.tsv" {
		t.Fatalf("unexpected inputFile: %q", task.InputFile)
	}
	if task.Metadata["stmtdate"] != "2026-03" {
		t.Fatalf("expected stmtdate metadata %q, got %q", "2026-03", task.Metadata["stmtdate"])
	}
}
