package check

import (
	"context"
	"errors"
	"testing"
)

func TestRedisCheck_Pass(t *testing.T) {
	c := &RedisCheck{URL: "redis://localhost:6379", pinger: func(_ context.Context, _ string) error {
		return nil
	}}
	result := c.Run(context.Background())
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %v: %s", result.Status, result.Message)
	}
}

func TestRedisCheck_Fail(t *testing.T) {
	c := &RedisCheck{URL: "redis://localhost:6379", pinger: func(_ context.Context, _ string) error {
		return errors.New("connection refused")
	}}
	result := c.Run(context.Background())
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %v", result.Status)
	}
}
