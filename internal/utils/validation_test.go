package utils

import (
	"testing"
)

type validationSample struct {
	Phone string `json:"phone" validate:"required,phone"`
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"min=2"`
}

func TestValidationErrors(t *testing.T) {
	sample := validationSample{
		Phone: "123",
		Email: "bad",
		Name:  "a",
	}

	err := Validate.Struct(sample)
	if err == nil {
		t.Fatalf("expected validation error")
	}

	details := ValidationErrors(err)
	if details == nil {
		t.Fatalf("expected validation details")
	}

	if details["phone"] != "invalid phone" {
		t.Fatalf("expected phone invalid, got %q", details["phone"])
	}
	if details["email"] != "invalid email" {
		t.Fatalf("expected email invalid, got %q", details["email"])
	}
	if details["name"] != "minimum" {
		t.Fatalf("expected name minimum, got %q", details["name"])
	}
}

func TestValidationErrorsUnknown(t *testing.T) {
	if details := ValidationErrors(nil); details != nil {
		t.Fatalf("expected nil details, got %#v", details)
	}
}
