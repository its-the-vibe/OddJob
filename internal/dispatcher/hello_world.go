package dispatcher

import (
	"fmt"
	"strings"

	"github.com/its-the-vibe/OddJob/internal/config"
)

const helloTaskName = "hello:world"

type HelloWorldTransformer struct{}

func NewHelloWorldTransformer() *HelloWorldTransformer {
	return &HelloWorldTransformer{}
}

func (h *HelloWorldTransformer) TaskName() string {
	return helloTaskName
}

func (h *HelloWorldTransformer) ToPoppit(task TaskMessage, cfg config.PoppitConfig) (PoppitMessage, error) {
	if task.TaskName != helloTaskName {
		return PoppitMessage{}, fmt.Errorf("unsupported task for hello transformer: %s", task.TaskName)
	}

	message := "hello world"
	if name := strings.TrimSpace(task.Metadata["name"]); name != "" {
		message = fmt.Sprintf("hello %s", name)
	}
	if input := strings.TrimSpace(task.InputFile); input != "" {
		message = fmt.Sprintf("%s from %s", message, input)
	}

	return PoppitMessage{
		Repo:     cfg.Repo,
		Branch:   cfg.Branch,
		Type:     cfg.Type,
		Dir:      cfg.Dir,
		Commands: []string{fmt.Sprintf("echo %q", message)},
		Metadata: map[string]string{
			"taskName": helloTaskName,
		},
	}, nil
}

func (h *HelloWorldTransformer) FromPoppit(output PoppitOutput) (*TaskMessage, bool, error) {
	if output.Metadata["taskName"] != helloTaskName {
		return nil, false, nil
	}
	if output.StatusCode != 0 {
		return nil, false, nil
	}

	nextTask := strings.TrimSpace(output.Metadata["nextTask"])
	if nextTask == "" {
		return nil, false, nil
	}

	return &TaskMessage{
		TaskName:  nextTask,
		InputFile: strings.TrimSpace(output.Output),
		Metadata: map[string]string{
			"sourceTask": helloTaskName,
		},
	}, true, nil
}
