package dispatcher

import (
	"testing"

	"github.com/its-the-vibe/OddJob/internal/config"
)

func TestSantanderPdftoppmToPoppit(t *testing.T) {
	transformer := NewSantanderPdftoppmTransformer()
	cfg := config.PoppitConfig{
		Repo:   "its-the-vibe/OddJob",
		Branch: "refs/heads/main",
		Dir:    "/workspace",
		Type:   "odd:job",
	}

	msg, err := transformer.ToPoppit(TaskMessage{
		TaskName:  santanderPdftoppmTaskName,
		InputFile: "/workspace/incoming/File-2026-03.pdf",
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
	if msg.Commands[0] != `pdftoppm -png -r 300 "File-2026-03.pdf" -f 2 "File-2026-03"` {
		t.Fatalf("unexpected command: %q", msg.Commands[0])
	}
	if msg.Metadata["taskName"] != santanderPdftoppmTaskName {
		t.Fatalf("expected metadata taskName %q, got %q", santanderPdftoppmTaskName, msg.Metadata["taskName"])
	}
	if msg.Metadata["stmtdate"] != "2026-03" {
		t.Fatalf("expected metadata stmtdate %q, got %q", "2026-03", msg.Metadata["stmtdate"])
	}
}

func TestSantanderPdftoppmToPoppitReturnsErrorWhenInputFileEmpty(t *testing.T) {
	transformer := NewSantanderPdftoppmTransformer()
	_, err := transformer.ToPoppit(TaskMessage{
		TaskName: santanderPdftoppmTaskName,
	}, config.PoppitConfig{})
	if err == nil {
		t.Fatalf("expected error when inputFile is empty")
	}
}

func TestSantanderPdftoppmToPoppitSupportsInputFileWithSpaces(t *testing.T) {
	transformer := NewSantanderPdftoppmTransformer()

	msg, err := transformer.ToPoppit(TaskMessage{
		TaskName:  santanderPdftoppmTaskName,
		InputFile: "/workspace/incoming/My File-2026-03.pdf",
	}, config.PoppitConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Dir != "/workspace/incoming" {
		t.Fatalf("unexpected dir: %q", msg.Dir)
	}
	if msg.Commands[0] != `pdftoppm -png -r 300 "My File-2026-03.pdf" -f 2 "My File-2026-03"` {
		t.Fatalf("unexpected command: %q", msg.Commands[0])
	}
}
