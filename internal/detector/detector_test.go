package detector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetect_PackageManager_npm(t *testing.T) {
	dir := t.TempDir()
	touch(t, filepath.Join(dir, "package.json"))
	// no lockfile → defaults to npm
	stack := Detect(dir)
	if !stack.Node {
		t.Fatal("expected Node to be detected")
	}
	if stack.PackageManager != "npm" {
		t.Errorf("expected PackageManager=npm, got %q", stack.PackageManager)
	}
}

func TestDetect_PackageManager_pnpm(t *testing.T) {
	dir := t.TempDir()
	touch(t, filepath.Join(dir, "package.json"))
	touch(t, filepath.Join(dir, "pnpm-lock.yaml"))
	stack := Detect(dir)
	if stack.PackageManager != "pnpm" {
		t.Errorf("expected PackageManager=pnpm, got %q", stack.PackageManager)
	}
}

func TestDetect_PackageManager_yarn(t *testing.T) {
	dir := t.TempDir()
	touch(t, filepath.Join(dir, "package.json"))
	touch(t, filepath.Join(dir, "yarn.lock"))
	stack := Detect(dir)
	if stack.PackageManager != "yarn" {
		t.Errorf("expected PackageManager=yarn, got %q", stack.PackageManager)
	}
}

func TestDetect_PackageManager_pnpm_takes_priority_over_yarn(t *testing.T) {
	// If somehow both lockfiles exist, pnpm wins (checked first).
	dir := t.TempDir()
	touch(t, filepath.Join(dir, "package.json"))
	touch(t, filepath.Join(dir, "pnpm-lock.yaml"))
	touch(t, filepath.Join(dir, "yarn.lock"))
	stack := Detect(dir)
	if stack.PackageManager != "pnpm" {
		t.Errorf("expected PackageManager=pnpm when both lockfiles exist, got %q", stack.PackageManager)
	}
}

func TestDetect_PackageManager_empty_when_not_node(t *testing.T) {
	dir := t.TempDir()
	// no package.json
	stack := Detect(dir)
	if stack.Node {
		t.Fatal("did not expect Node to be detected")
	}
	if stack.PackageManager != "" {
		t.Errorf("expected PackageManager to be empty for non-Node project, got %q", stack.PackageManager)
	}
}

func touch(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatalf("touch %s: %v", path, err)
	}
}