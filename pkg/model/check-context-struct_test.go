package model

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPassword_ValueOnly(t *testing.T) {
	cc := CheckContext{PasswordValue: "secret123"}
	got := cc.Password()
	if got != "secret123" {
		t.Errorf("Password() = %q, want %q", got, "secret123")
	}
}

func TestPassword_EmptyValue(t *testing.T) {
	cc := CheckContext{PasswordValue: ""}
	got := cc.Password()
	if got != "" {
		t.Errorf("Password() = %q, want empty string", got)
	}
}

func TestPassword_FileTakesPrecedence(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "pass.txt")
	if err := os.WriteFile(f, []byte("file-password\n"), 0644); err != nil {
		t.Fatalf("failed to write password file: %v", err)
	}

	cc := CheckContext{
		PasswordValue: "value-password",
		PasswordFile:  f,
	}
	got := cc.Password()
	if got != "file-password" {
		t.Errorf("Password() = %q, want %q (file should take precedence)", got, "file-password")
	}
}

func TestPassword_FileTrimsWhitespace(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "pass.txt")
	if err := os.WriteFile(f, []byte("  my-secret  \n"), 0644); err != nil {
		t.Fatalf("failed to write password file: %v", err)
	}

	cc := CheckContext{PasswordFile: f}
	got := cc.Password()
	if got != "my-secret" {
		t.Errorf("Password() = %q, want %q (whitespace should be trimmed)", got, "my-secret")
	}
}

func TestPassword_NonexistentFile(t *testing.T) {
	cc := CheckContext{
		PasswordValue: "fallback",
		PasswordFile:  "/nonexistent/path/pass.txt",
	}
	got := cc.Password()
	// When file doesn't exist, os.ReadFile fails and returns empty string
	// (the function logs the error but does not fall back to PasswordValue)
	if got != "" {
		t.Errorf("Password() = %q, want empty string (file read failure returns empty)", got)
	}
}

func TestPassword_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "pass.txt")
	if err := os.WriteFile(f, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write password file: %v", err)
	}

	cc := CheckContext{PasswordFile: f}
	got := cc.Password()
	if got != "" {
		t.Errorf("Password() = %q, want empty string", got)
	}
}

func TestValidate_PasswordFileExists(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "pass.txt")
	if err := os.WriteFile(f, []byte("secret"), 0644); err != nil {
		t.Fatalf("failed to write password file: %v", err)
	}

	cc := CheckContext{PasswordFile: f}
	if err := cc.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

func TestValidate_PasswordFileDoesNotExist(t *testing.T) {
	cc := CheckContext{PasswordFile: "/nonexistent/path/pass.txt"}
	if err := cc.Validate(); err == nil {
		t.Error("Validate() should return error for nonexistent password file")
	}
}

func TestValidate_NoPasswordFile(t *testing.T) {
	cc := CheckContext{}
	if err := cc.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}
