package dispatcher

import (
	"fmt"

	"github.com/its-the-vibe/OddJob/internal/config"
)

type Transformer interface {
	TaskName() string
	ToPoppit(task TaskMessage, cfg config.PoppitConfig) (PoppitMessage, error)
	FromPoppit(output PoppitOutput) (*TaskMessage, bool, error)
}

type Registry struct {
	byTask map[string]Transformer
	all    []Transformer
}

func NewRegistry(transformers ...Transformer) *Registry {
	registry := &Registry{
		byTask: make(map[string]Transformer, len(transformers)),
		all:    transformers,
	}
	for _, transformer := range transformers {
		registry.byTask[transformer.TaskName()] = transformer
	}
	return registry
}

func (r *Registry) ToPoppit(task TaskMessage, cfg config.PoppitConfig) (PoppitMessage, error) {
	transformer, ok := r.byTask[task.TaskName]
	if !ok {
		return PoppitMessage{}, fmt.Errorf("no transformer registered for task %q", task.TaskName)
	}
	return transformer.ToPoppit(task, cfg)
}

func (r *Registry) FromPoppit(output PoppitOutput) (*TaskMessage, bool, error) {
	for _, transformer := range r.all {
		task, ok, err := transformer.FromPoppit(output)
		if err != nil {
			return nil, false, err
		}
		if ok {
			return task, true, nil
		}
	}
	return nil, false, nil
}
