package dispatcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/its-the-vibe/OddJob/internal/config"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	cfg      config.Config
	client   *redis.Client
	registry *Registry
	logger   *log.Logger
}

func NewService(cfg config.Config, client *redis.Client, registry *Registry, logger *log.Logger) *Service {
	return &Service{cfg: cfg, client: client, registry: registry, logger: logger}
}

func (s *Service) Run(ctx context.Context) error {
	errCh := make(chan error, 2)

	go func() {
		errCh <- s.consumeTasks(ctx)
	}()

	go func() {
		errCh <- s.consumePoppitOutput(ctx)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			if err == nil || errors.Is(err, context.Canceled) {
				continue
			}
			return err
		}
	}
}

func (s *Service) consumeTasks(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		payload, err := s.client.LPop(ctx, s.cfg.Queues.TaskQueue).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				time.Sleep(s.cfg.PollInterval())
				continue
			}
			return fmt.Errorf("lpop task queue: %w", err)
		}

		var task TaskMessage
		if err := json.Unmarshal([]byte(payload), &task); err != nil {
			s.logger.Printf("discarding malformed task payload: %v", err)
			continue
		}

		message, err := s.registry.ToPoppit(task, s.cfg.Poppit)
		if err != nil {
			s.logger.Printf("skip task %q: %v", task.TaskName, err)
			continue
		}

		encoded, err := json.Marshal(message)
		if err != nil {
			s.logger.Printf("encode poppit message failed: %v", err)
			continue
		}

		if err := s.client.RPush(ctx, s.cfg.Queues.PoppitQueue, encoded).Err(); err != nil {
			return fmt.Errorf("rpush poppit queue: %w", err)
		}
	}
}

func (s *Service) consumePoppitOutput(ctx context.Context) error {
	pubSub := s.client.Subscribe(ctx, s.cfg.Poppit.OutputChannel)
	defer pubSub.Close()

	if _, err := pubSub.Receive(ctx); err != nil {
		return fmt.Errorf("subscribe poppit output: %w", err)
	}

	channel := pubSub.Channel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case message, ok := <-channel:
			if !ok {
				return fmt.Errorf("poppit output channel closed")
			}

			var output PoppitOutput
			if err := json.Unmarshal([]byte(message.Payload), &output); err != nil {
				s.logger.Printf("discarding malformed poppit output: %v", err)
				continue
			}
			if output.Type != s.cfg.Poppit.Type {
				continue
			}

			task, ok, err := s.registry.FromPoppit(output)
			if err != nil {
				s.logger.Printf("transform poppit output failed: %v", err)
				continue
			}
			if !ok {
				continue
			}

			encoded, err := json.Marshal(task)
			if err != nil {
				s.logger.Printf("encode chained task failed: %v", err)
				continue
			}
			if err := s.client.RPush(ctx, s.cfg.Queues.TaskQueue, encoded).Err(); err != nil {
				return fmt.Errorf("rpush task queue: %w", err)
			}
		}
	}
}
