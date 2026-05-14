package dispatcher

import (
	"testing"

	"github.com/its-the-vibe/OddJob/internal/config"
)

func TestSantanderStmt2redisToPoppit(t *testing.T) {
	transformer := NewSantanderStmt2redisTransformer()
	cfg := config.PoppitConfig{
		Repo:   "its-the-vibe/OddJob",
		Branch: "refs/heads/main",
		Dir:    "/workspace",
		Type:   "odd:job",
	}

	msg, err := transformer.ToPoppit(TaskMessage{
		TaskName:  santanderStmt2redisTaskName,
		InputFile: "/workspace/incoming/File-2026-03-2.tsv",
		Metadata: map[string]string{
			"stmtdate": "2026-03",
		},
	}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg.Dir != santanderStmt2redisDir {
		t.Fatalf("expected dir %q, got %q", santanderStmt2redisDir, msg.Dir)
	}
	if len(msg.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(msg.Commands))
	}
	if msg.Commands[0] != `${stmt2redis} -f "/workspace/incoming/File-2026-03-2.tsv" -t santander` {
		t.Fatalf("unexpected command: %q", msg.Commands[0])
	}
	if msg.Metadata["taskName"] != santanderStmt2redisTaskName {
		t.Fatalf("expected metadata taskName %q, got %q", santanderStmt2redisTaskName, msg.Metadata["taskName"])
	}
	if msg.Metadata["stmtdate"] != "2026-03" {
		t.Fatalf("expected metadata stmtdate %q, got %q", "2026-03", msg.Metadata["stmtdate"])
	}
}

func TestSantanderStmt2redisToPoppitReturnsErrorWhenInputFileEmpty(t *testing.T) {
	transformer := NewSantanderStmt2redisTransformer()
	_, err := transformer.ToPoppit(TaskMessage{
		TaskName: santanderStmt2redisTaskName,
	}, config.PoppitConfig{})
	if err == nil {
		t.Fatalf("expected error when inputFile is empty")
	}
}
