package dispatcher

import (
	"testing"

	"github.com/its-the-vibe/OddJob/internal/config"
)

func TestSantanderStmtdateToPoppit(t *testing.T) {
	transformer := NewSantanderStmtdateTransformer()
	cfg := config.PoppitConfig{
		Repo:   "its-the-vibe/OddJob",
		Branch: "refs/heads/main",
		Dir:    "/workspace",
		Type:   "odd:job",
	}

	msg, err := transformer.ToPoppit(TaskMessage{
		TaskName:  santanderStmtdateTaskName,
		InputFile: "statement.pdf",
		Metadata:  map[string]string{"source": "test"},
	}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(msg.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(msg.Commands))
	}
	if msg.Commands[0] != `stmtdate -rename "statement.pdf"` {
		t.Fatalf("unexpected command: %q", msg.Commands[0])
	}
	if msg.Metadata["source"] != "test" {
		t.Fatalf("expected metadata source %q, got %q", "test", msg.Metadata["source"])
	}
}

func TestSantanderStmtdateToPoppitRequiresInputFile(t *testing.T) {
	transformer := NewSantanderStmtdateTransformer()
	_, err := transformer.ToPoppit(TaskMessage{
		TaskName: santanderStmtdateTaskName,
	}, config.PoppitConfig{})
	if err == nil {
		t.Fatalf("expected error when inputFile is empty")
	}
}
