package utils

import "testing"

func TestMaskCitizenID(t *testing.T) {
	if MaskCitizenID("") != "" {
		t.Fatalf("expected empty mask")
	}
	if MaskCitizenID("1234") != "****" {
		t.Fatalf("expected full mask for short id")
	}
	if MaskCitizenID("1234567890123") != "*********0123" {
		t.Fatalf("unexpected mask")
	}
}
