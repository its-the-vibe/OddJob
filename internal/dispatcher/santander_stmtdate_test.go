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
	if msg.Metadata["taskName"] != santanderStmtdateTaskName {
		t.Fatalf("expected metadata taskName %q, got %q", santanderStmtdateTaskName, msg.Metadata["taskName"])
	}
}

func TestSantanderStmtdateToPoppitReturnsErrorWhenInputFileEmpty(t *testing.T) {
	transformer := NewSantanderStmtdateTransformer()
	_, err := transformer.ToPoppit(TaskMessage{
		TaskName: santanderStmtdateTaskName,
	}, config.PoppitConfig{})
	if err == nil {
		t.Fatalf("expected error when inputFile is empty")
	}
}

func TestSantanderStmtdateFromPoppitChainsToPdftoppm(t *testing.T) {
	transformer := NewSantanderStmtdateTransformer()

	task, ok, err := transformer.FromPoppit(PoppitOutput{
		StatusCode: 0,
		Output: `2026-03
renamed: /workspace/incoming/File-2026-03.pdf`,
		Metadata: map[string]string{
			"taskName": santanderStmtdateTaskName,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected chained task")
	}
	if task.TaskName != santanderPdftoppmTaskName {
		t.Fatalf("expected taskName %q, got %q", santanderPdftoppmTaskName, task.TaskName)
	}
	if task.InputFile != "/workspace/incoming/File-2026-03.pdf" {
		t.Fatalf("unexpected inputFile: %q", task.InputFile)
	}
	if task.Metadata["stmtdate"] != "2026-03" {
		t.Fatalf("expected stmtdate metadata %q, got %q", "2026-03", task.Metadata["stmtdate"])
	}
}

func TestSantanderStmtdateFromPoppitReturnsErrorOnMalformedOutput(t *testing.T) {
	transformer := NewSantanderStmtdateTransformer()

	task, ok, err := transformer.FromPoppit(PoppitOutput{
		StatusCode: 0,
		Output:     "2026-03",
		Metadata: map[string]string{
			"taskName": santanderStmtdateTaskName,
		},
	})
	if err == nil {
		t.Fatalf("expected error for malformed output")
	}
	if ok {
		t.Fatalf("expected no chained task")
	}
	if task != nil {
		t.Fatalf("expected no task")
	}
}

func TestSantanderStmtdateFromPoppitParsesRenamedPathWithSpaces(t *testing.T) {
	transformer := NewSantanderStmtdateTransformer()

	task, ok, err := transformer.FromPoppit(PoppitOutput{
		StatusCode: 0,
		Output: `2026-03
renamed: /workspace/incoming/My File-2026-03.pdf`,
		Metadata: map[string]string{
			"taskName": santanderStmtdateTaskName,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected chained task")
	}
	if task.InputFile != "/workspace/incoming/My File-2026-03.pdf" {
		t.Fatalf("unexpected inputFile: %q", task.InputFile)
	}
}
