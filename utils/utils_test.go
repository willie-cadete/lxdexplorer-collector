package utils

import "testing"

func TestCapitalizeWithLower(t *testing.T) {
	capitalized := Capitalize("test")
	if capitalized != "Test" {
		t.Errorf("Expected 'Test', got %s", capitalized)
	}
}

func TestCapitalizeWithUpper(t *testing.T) {
	capitalized := Capitalize("TEST")
	if capitalized != "Test" {
		t.Errorf("Expected 'Test', got %s", capitalized)
	}
}

func TestCapitalizeWithEmptyString(t *testing.T) {
	capitalized := Capitalize("")
	if capitalized != "" {
		t.Errorf("Expected '', got %s", capitalized)
	}
}
