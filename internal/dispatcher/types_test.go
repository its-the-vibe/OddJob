package dispatcher

import (
	"encoding/json"
	"testing"
)

func TestPoppitOutputMixedTypeMetadata(t *testing.T) {
	raw := `{"metadata":{"taskName":"x","count":3,"addToProject":true},"type":"t","command":"c","output":"","stderr":"","status_code":0}`
	var out PoppitOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := metadataString(out.Metadata, "count"); got != "3" {
		t.Fatalf("expected 3, got %q", got)
	}
	if got := metadataString(out.Metadata, "addToProject"); got != "true" {
		t.Fatalf("expected true, got %q", got)
	}
	if got := metadataString(out.Metadata, "missing"); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}
