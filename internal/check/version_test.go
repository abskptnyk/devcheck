package check

import (
	"context"
	"os"
	"testing"
)

func TestVersionLess(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"1.21.0", "1.22", true},
		{"1.22.0", "1.21", false},
		{"1.22.0", "1.22", false},
		{"20.0.0", "18", false},
		{"16.0.0", "18", true},
	}
	for _, tc := range cases {
		got := versionLess(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("versionLess(%q, %q) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestGoVersionCheck_Pass(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/go.mod", []byte("module example\n\ngo 1.1\n"), 0644)

	c := &GoVersionCheck{Dir: dir}
	result := c.Run(context.Background())
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %v: %s", result.Status, result.Message)
	}
}

func TestGoVersionCheck_Fail(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/go.mod", []byte("module example\n\ngo 9999.0\n"), 0644)

	c := &GoVersionCheck{Dir: dir}
	result := c.Run(context.Background())
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %v: %s", result.Status, result.Message)
	}
}

func TestGoVersionCheck_MissingGoMod(t *testing.T) {
	c := &GoVersionCheck{Dir: t.TempDir()}
	result := c.Run(context.Background())
	if result.Status != StatusSkipped {
		t.Errorf("expected skipped, got %v", result.Status)
	}
}

func TestNodeVersionCheck_NvmrcPass(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/.nvmrc", []byte("1\n"), 0644)

	c := &NodeVersionCheck{Dir: dir}
	result := c.Run(context.Background())
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %v: %s", result.Status, result.Message)
	}
}

func TestNodeVersionCheck_PackageJsonPass(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/package.json", []byte(`{"engines":{"node":">=1.0.0"}}`), 0644)

	c := &NodeVersionCheck{Dir: dir}
	result := c.Run(context.Background())
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %v: %s", result.Status, result.Message)
	}
}

func TestNodeVersionCheck_NoRequirement(t *testing.T) {
	c := &NodeVersionCheck{Dir: t.TempDir()}
	result := c.Run(context.Background())
	if result.Status != StatusSkipped {
		t.Errorf("expected skipped, got %v", result.Status)
	}
}
