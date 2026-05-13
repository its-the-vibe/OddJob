package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadUsesEnvPasswordAndDefaults(t *testing.T) {
	t.Setenv("REDIS_PASSWORD", "from-env")

	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	configJSON := `{
  "redis": {"addr": "localhost:6379", "db": 0},
  "queues": {"taskQueue": "oddjob:tasks", "poppitQueue": "poppit:in"},
  "poppit": {
    "repo": "its-the-vibe/OddJob",
    "branch": "refs/heads/main",
    "dir": "/workdir",
    "outputChannel": "poppit:out"
  },
  "poll": {}
}`
	if err := os.WriteFile(configPath, []byte(configJSON), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Redis.Password != "from-env" {
		t.Fatalf("expected env password override, got %q", cfg.Redis.Password)
	}
	if cfg.Poppit.Type != "odd:job" {
		t.Fatalf("expected default type odd:job, got %q", cfg.Poppit.Type)
	}
	if cfg.Poll.IntervalSeconds != 1 {
		t.Fatalf("expected default poll interval 1, got %d", cfg.Poll.IntervalSeconds)
	}
}

func TestLoadResolvesAliases(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	configJSON := `{
  "redis": {"addr": "${aliases.redisAddr}", "db": 0},
  "queues": {"taskQueue": "${aliases.taskQueue}", "poppitQueue": "poppit:in"},
  "poppit": {
    "repo": "its-the-vibe/OddJob",
    "branch": "refs/heads/main",
    "dir": "${aliases.incomingDir}",
    "type": "${aliases.messageType}",
    "outputChannel": "poppit:out"
  },
  "poll": {},
  "aliases": {
    "baseDir": "/workspace",
    "incomingDir": "${aliases.baseDir}/incoming",
    "redisAddr": "localhost:6379",
    "taskQueue": "oddjob:tasks",
    "messageType": "odd:job"
  }
}`
	if err := os.WriteFile(configPath, []byte(configJSON), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Poppit.Dir != "/workspace/incoming" {
		t.Fatalf("expected poppit.dir to resolve aliases, got %q", cfg.Poppit.Dir)
	}
	if cfg.Queues.TaskQueue != "oddjob:tasks" {
		t.Fatalf("expected queues.taskQueue to resolve aliases, got %q", cfg.Queues.TaskQueue)
	}
	if cfg.Redis.Addr != "localhost:6379" {
		t.Fatalf("expected redis.addr to resolve aliases, got %q", cfg.Redis.Addr)
	}
	if cfg.Aliases["incomingDir"] != "/workspace/incoming" {
		t.Fatalf("expected aliases.incomingDir to resolve nested aliases, got %q", cfg.Aliases["incomingDir"])
	}
}

func TestLoadReturnsErrorForUndefinedAlias(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	configJSON := `{
  "redis": {"addr": "localhost:6379", "db": 0},
  "queues": {"taskQueue": "oddjob:tasks", "poppitQueue": "poppit:in"},
  "poppit": {
    "repo": "its-the-vibe/OddJob",
    "branch": "refs/heads/main",
    "dir": "${aliases.missing}",
    "outputChannel": "poppit:out"
  },
  "poll": {},
  "aliases": {
    "baseDir": "/workspace"
  }
}`
	if err := os.WriteFile(configPath, []byte(configJSON), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	if _, err := Load(configPath); err == nil {
		t.Fatalf("expected load config to fail for undefined alias")
	}
}
