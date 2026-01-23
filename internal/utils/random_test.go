package utils

import (
	"regexp"
	"testing"
)

func TestRandomDigits(t *testing.T) {
	value, err := RandomDigits(6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(value) != 6 {
		t.Fatalf("expected length 6, got %d", len(value))
	}
	if !regexp.MustCompile(`^[0-9]+$`).MatchString(value) {
		t.Fatalf("expected digits only, got %q", value)
	}
}

func TestRandomDigitsInvalid(t *testing.T) {
	if _, err := RandomDigits(0); err == nil {
		t.Fatalf("expected error for invalid length")
	}
}

func TestRandomRefCode(t *testing.T) {
	value, err := RandomRefCode(8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(value) != 8 {
		t.Fatalf("expected length 8, got %d", len(value))
	}
	if !regexp.MustCompile(`^[0-9a-f]+$`).MatchString(value) {
		t.Fatalf("expected hex, got %q", value)
	}
}

func TestRandomRefCodeInvalid(t *testing.T) {
	if _, err := RandomRefCode(0); err == nil {
		t.Fatalf("expected error for invalid length")
	}
}
