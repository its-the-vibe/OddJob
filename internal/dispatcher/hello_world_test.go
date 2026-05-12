package dispatcher

import (
	"testing"

	"github.com/its-the-vibe/OddJob/internal/config"
)

func TestHelloWorldToPoppit(t *testing.T) {
	transformer := NewHelloWorldTransformer()
	cfg := config.PoppitConfig{
		Repo:   "its-the-vibe/OddJob",
		Branch: "refs/heads/main",
		Dir:    "/workspace",
		Type:   "odd:job",
	}

	msg, err := transformer.ToPoppit(TaskMessage{
		TaskName: helloTaskName,
		Metadata: map[string]string{"name": "OddJob"},
	}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(msg.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(msg.Commands))
	}
	// if msg.Metadata["taskName"] != helloTaskName {
	// 	t.Fatalf("expected taskName metadata %q, got %q", helloTaskName, msg.Metadata["taskName"])
	// }
	if msg.Metadata["name"] != "OddJob" {
		t.Fatalf("expected name metadata %q, got %q", "OddJob", msg.Metadata["name"])
	}

}

func TestHelloWorldFromPoppitChainsWhenRequested(t *testing.T) {
	transformer := NewHelloWorldTransformer()

	task, ok, err := transformer.FromPoppit(PoppitOutput{
		Type:       "odd:job",
		StatusCode: 0,
		Output:     "next-input",
		Metadata: map[string]string{
			"taskName": helloTaskName,
			"nextTask": "downstream:task",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected chaining to be enabled")
	}
	if task.TaskName != "downstream:task" {
		t.Fatalf("expected downstream task name, got %q", task.TaskName)
	}
}
