package id_generator

import (
	"fmt"
	"strings"
	"testing"
)

func TestIDGenerator_DefaultConfig(t *testing.T) {
	gen := NewIDGenerator(GetDefaultConfig())
	id := gen.Generate()
	if id == "" {
		t.Error("ID should not be empty")
	}
	parts := strings.Split(id, "-")
	if len(parts) != 3 {
		t.Errorf("Default ID should have 3 parts, got %d", len(parts))
	}
}

func TestIDGenerator_CustomConfig(t *testing.T) {
	sep := "_"
	customModifier := func(timestamp int64, microSecond int64, randSegment string) []string {
		return []string{
			"CUSTOM",
			fmt.Sprintf("%d", timestamp),
			randSegment,
		}
	}
	cfg := Config{
		Separator: &sep,
		Modifier:  customModifier,
	}
	gen := NewIDGenerator(cfg)
	id := gen.Generate()
	if id == "" {
		t.Error("Custom ID should not be empty")
	}
	parts := strings.Split(id, sep)
	if len(parts) != 3 {
		t.Errorf("Custom ID should have 3 parts, got %d", len(parts))
	}
	if parts[0] != "CUSTOM" {
		t.Errorf("Custom ID first part should be 'CUSTOM', got %s", parts[0])
	}
}
