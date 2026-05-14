package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	Redis   RedisConfig       `json:"redis"`
	Queues  QueueConfig       `json:"queues"`
	Poppit  PoppitConfig      `json:"poppit"`
	Poll    PollConfig        `json:"poll"`
	Aliases map[string]string `json:"aliases"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type QueueConfig struct {
	TaskQueue   string `json:"taskQueue"`
	PoppitQueue string `json:"poppitQueue"`
}

type PoppitConfig struct {
	Repo          string `json:"repo"`
	Branch        string `json:"branch"`
	Dir           string `json:"dir"`
	Type          string `json:"type"`
	OutputChannel string `json:"outputChannel"`
}

type PollConfig struct {
	IntervalSeconds int `json:"intervalSeconds"`
}

func Load(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(content, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config json: %w", err)
	}

	if envPassword := os.Getenv("REDIS_PASSWORD"); envPassword != "" {
		cfg.Redis.Password = envPassword
	}

	if cfg.Poll.IntervalSeconds <= 0 {
		cfg.Poll.IntervalSeconds = 1
	}
	if cfg.Poppit.Type == "" {
		cfg.Poppit.Type = "odd:job"
	}

	cfg.Poppit.Repo = cfg.ResolveAliases(cfg.Poppit.Repo)
	cfg.Poppit.Branch = cfg.ResolveAliases(cfg.Poppit.Branch)
	cfg.Poppit.Dir = cfg.ResolveAliases(cfg.Poppit.Dir)
	cfg.Poppit.Type = cfg.ResolveAliases(cfg.Poppit.Type)
	cfg.Poppit.OutputChannel = cfg.ResolveAliases(cfg.Poppit.OutputChannel)

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// ResolveAliases replaces ${key} placeholders in s with the corresponding alias value.
func (c Config) ResolveAliases(s string) string {
	if len(c.Aliases) == 0 {
		return s
	}
	pairs := make([]string, 0, len(c.Aliases)*2)
	for key, value := range c.Aliases {
		pairs = append(pairs, "${"+key+"}", value)
	}
	return strings.NewReplacer(pairs...).Replace(s)
}

func (c Config) PollInterval() time.Duration {
	return time.Duration(c.Poll.IntervalSeconds) * time.Second
}

func (c Config) validate() error {
	switch {
	case c.Redis.Addr == "":
		return fmt.Errorf("redis.addr is required")
	case c.Queues.TaskQueue == "":
		return fmt.Errorf("queues.taskQueue is required")
	case c.Queues.PoppitQueue == "":
		return fmt.Errorf("queues.poppitQueue is required")
	case c.Poppit.OutputChannel == "":
		return fmt.Errorf("poppit.outputChannel is required")
	case c.Poppit.Repo == "":
		return fmt.Errorf("poppit.repo is required")
	case c.Poppit.Branch == "":
		return fmt.Errorf("poppit.branch is required")
	case c.Poppit.Dir == "":
		return fmt.Errorf("poppit.dir is required")
	default:
		return nil
	}
}
