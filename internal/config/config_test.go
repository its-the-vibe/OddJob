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

func TestLoadResolvesAliasesInPoppitConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	configJSON := `{
  "redis": {"addr": "localhost:6379", "db": 0},
  "queues": {"taskQueue": "oddjob:tasks", "poppitQueue": "poppit:in"},
  "poppit": {
    "repo": "its-the-vibe/OddJob",
    "branch": "refs/heads/main",
    "dir": "${basedir}",
    "outputChannel": "poppit:out"
  },
  "poll": {},
  "aliases": {
    "basedir": "/data/workspace"
  }
}`
	if err := os.WriteFile(configPath, []byte(configJSON), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Poppit.Dir != "/data/workspace" {
		t.Fatalf("expected alias-resolved dir %q, got %q", "/data/workspace", cfg.Poppit.Dir)
	}
}

func TestResolveAliasesSubstitutesPlaceholders(t *testing.T) {
	cfg := Config{
		Aliases: map[string]string{
			"basedir": "/data/workspace",
			"cmd":     "/usr/local/bin/tool",
		},
	}

	tests := []struct {
		input string
		want  string
	}{
		{"${basedir}/file.txt", "/data/workspace/file.txt"},
		{"${cmd} --flag", "/usr/local/bin/tool --flag"},
		{"no placeholder", "no placeholder"},
		{"${basedir}/${basedir}", "/data/workspace//data/workspace"},
		{"${unknown}", "${unknown}"},
	}

	for _, tt := range tests {
		got := cfg.ResolveAliases(tt.input)
		if got != tt.want {
			t.Errorf("ResolveAliases(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
