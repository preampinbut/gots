package main

import (
	"strings"
	"testing"
)

func TestTranspileRecursive(t *testing.T) {
	got, err := TranspileFile("input/input1.go")
	if err != nil {
		t.Fatalf("failed to transpile: %v", err)
	}

	expected := `
interface Address {
  street: string;
  city: string;
}

interface Person {
  name: string;
  age: number;
  address: Address;
  friends: Person[];
}
`

	// Normalize whitespace
	normalize := func(s string) string {
		return strings.TrimSpace(strings.ReplaceAll(s, "\r\n", "\n"))
	}

	if normalize(got) != normalize(expected) {
		t.Errorf("transpiled output mismatch\nGot:\n%s\n\nExpected:\n%s", got, expected)
	}
}

func TestTranspileOmitEmpry(t *testing.T) {
	got, err := TranspileFile("input/input2.go")
	if err != nil {
		t.Fatalf("failed to transpile: %v", err)
	}

	expected := `
interface PersonalInfo {
  name: string;
  tel?: string;
  email?: string;
}
`

	// Normalize whitespace
	normalize := func(s string) string {
		return strings.TrimSpace(strings.ReplaceAll(s, "\r\n", "\n"))
	}

	if normalize(got) != normalize(expected) {
		t.Errorf("transpiled output mismatch\nGot:\n%s\n\nExpected:\n%s", got, expected)
	}
}

func TestTranspileImport(t *testing.T) {
	got, err := TranspileFile("input/input3.go")
	if err != nil {
		t.Fatalf("failed to transpile with import: %v", err)
	}

	expected := `
interface Pet {
  name: string;
}

interface PetOwnership {
  name: string;
  pet: Pet[];
}
`

	normalize := func(s string) string {
		return strings.TrimSpace(strings.ReplaceAll(s, "\r\n", "\n"))
	}

	if normalize(got) != normalize(expected) {
		t.Errorf("transpiled output mismatch\nGot:\n%s\n\nExpected:\n%s", got, expected)
	}
}

func TestTranspileAliases(t *testing.T) {
	got, err := TranspileFile("input/input4.go")
	if err != nil {
		t.Fatalf("failed to transpile with import: %v", err)
	}

	expected := `
interface Card {
  id: number;
  email: string;
}
`

	normalize := func(s string) string {
		return strings.TrimSpace(strings.ReplaceAll(s, "\r\n", "\n"))
	}

	if normalize(got) != normalize(expected) {
		t.Errorf("transpiled output mismatch\nGot:\n%s\n\nExpected:\n%s", got, expected)
	}
}
