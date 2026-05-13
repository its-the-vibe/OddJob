package config

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"
)

var aliasPattern = regexp.MustCompile(`\$\{aliases\.([^}]+)\}`)

type Config struct {
	Redis   RedisConfig       `json:"redis"`
	Queues  QueueConfig       `json:"queues"`
	Poppit  PoppitConfig      `json:"poppit"`
	Poll    PollConfig        `json:"poll"`
	Aliases map[string]string `json:"aliases,omitempty"`
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

	if err := cfg.resolveAliases(); err != nil {
		return Config{}, err
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
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

func (c *Config) resolveAliases() error {
	if len(c.Aliases) == 0 {
		return nil
	}

	resolved := make(map[string]string, len(c.Aliases))
	resolving := make(map[string]bool, len(c.Aliases))

	var resolveAlias func(name string) (string, error)
	resolveAlias = func(name string) (string, error) {
		if value, ok := resolved[name]; ok {
			return value, nil
		}
		value, ok := c.Aliases[name]
		if !ok {
			return "", fmt.Errorf("aliases.%s is not defined", name)
		}
		if resolving[name] {
			return "", fmt.Errorf("aliases.%s contains a circular reference", name)
		}

		resolving[name] = true
		expanded, err := expandAliasesInString(value, resolveAlias)
		resolving[name] = false
		if err != nil {
			return "", err
		}

		resolved[name] = expanded
		return expanded, nil
	}

	for key := range c.Aliases {
		if _, err := resolveAlias(key); err != nil {
			return err
		}
	}

	resolveFromMap := func(name string) (string, error) {
		value, ok := resolved[name]
		if !ok {
			return "", fmt.Errorf("aliases.%s is not defined", name)
		}
		return value, nil
	}

	var err error
	if c.Redis.Addr, err = expandAliasesInString(c.Redis.Addr, resolveFromMap); err != nil {
		return err
	}
	if c.Redis.Username, err = expandAliasesInString(c.Redis.Username, resolveFromMap); err != nil {
		return err
	}
	if c.Redis.Password, err = expandAliasesInString(c.Redis.Password, resolveFromMap); err != nil {
		return err
	}
	if c.Queues.TaskQueue, err = expandAliasesInString(c.Queues.TaskQueue, resolveFromMap); err != nil {
		return err
	}
	if c.Queues.PoppitQueue, err = expandAliasesInString(c.Queues.PoppitQueue, resolveFromMap); err != nil {
		return err
	}
	if c.Poppit.Repo, err = expandAliasesInString(c.Poppit.Repo, resolveFromMap); err != nil {
		return err
	}
	if c.Poppit.Branch, err = expandAliasesInString(c.Poppit.Branch, resolveFromMap); err != nil {
		return err
	}
	if c.Poppit.Dir, err = expandAliasesInString(c.Poppit.Dir, resolveFromMap); err != nil {
		return err
	}
	if c.Poppit.Type, err = expandAliasesInString(c.Poppit.Type, resolveFromMap); err != nil {
		return err
	}
	if c.Poppit.OutputChannel, err = expandAliasesInString(c.Poppit.OutputChannel, resolveFromMap); err != nil {
		return err
	}

	c.Aliases = resolved
	return nil
}

func expandAliasesInString(value string, resolveAlias func(name string) (string, error)) (string, error) {
	matches := aliasPattern.FindAllStringSubmatchIndex(value, -1)
	if len(matches) == 0 {
		return value, nil
	}

	result := make([]byte, 0, len(value))
	last := 0
	for _, match := range matches {
		result = append(result, value[last:match[0]]...)

		aliasName := value[match[2]:match[3]]
		resolvedAlias, err := resolveAlias(aliasName)
		if err != nil {
			return "", err
		}
		result = append(result, resolvedAlias...)
		last = match[1]
	}
	result = append(result, value[last:]...)
	return string(result), nil
}
