package checks

import (
	"context"
	"testing"

	"github.com/vidya381/devcheck/internal/check"
)

func TestBinaryCheck_Pass(t *testing.T) {
	c := &BinaryCheck{Binary: "go"}
	result := c.Run(context.Background())
	if result.Status != check.StatusPass {
		t.Errorf("expected pass, got %v: %s", result.Status, result.Message)
	}
}

func TestBinaryCheck_Fail(t *testing.T) {
	c := &BinaryCheck{Binary: "definitelynotabinary12345"}
	result := c.Run(context.Background())
	if result.Status != check.StatusFail {
		t.Errorf("expected fail, got %v", result.Status)
	}
}
