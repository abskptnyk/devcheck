package check

import (
	"context"
	"errors"
	"testing"
)

func TestPortCheck_Pass(t *testing.T) {
	c := &PortCheck{Service: "PostgreSQL", Port: "5432", dialer: func(_ string) error {
		return nil
	}}
	result := c.Run(context.Background())
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %v: %s", result.Status, result.Message)
	}
}

func TestPortCheck_Fail(t *testing.T) {
	c := &PortCheck{Service: "PostgreSQL", Port: "5432", dialer: func(_ string) error {
		return errors.New("connection refused")
	}}
	result := c.Run(context.Background())
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %v", result.Status)
	}
}

func TestPortCheck_MessageContainsPort(t *testing.T) {
	c := &PortCheck{Service: "Redis", Port: "6379", dialer: func(_ string) error {
		return errors.New("connection refused")
	}}
	result := c.Run(context.Background())
	if result.Message != "nothing listening on port 6379" {
		t.Errorf("unexpected message: %s", result.Message)
	}
}

func TestPortFromURL_StandardURL(t *testing.T) {
	cases := []struct {
		url          string
		defaultPort  string
		expectedPort string
	}{
		{"postgres://user:pass@localhost:5555/db", "5432", "5555"},
		{"redis://localhost:6380", "6379", "6380"},
		{"mongodb://localhost:27018", "27017", "27018"},
		{"", "5432", "5432"},
		{"postgres://localhost/db", "5432", "5432"}, // no port in URL → default
	}
	for _, tc := range cases {
		got := portFromURL(tc.url, tc.defaultPort)
		if got != tc.expectedPort {
			t.Errorf("portFromURL(%q, %q) = %q, want %q", tc.url, tc.defaultPort, got, tc.expectedPort)
		}
	}
}

func TestPortFromURL_MySQLDSN(t *testing.T) {
	got := portFromURL("user:pass@tcp(localhost:3307)/db", "3306")
	if got != "3307" {
		t.Errorf("expected 3307, got %s", got)
	}
}

func TestPortFromURL_MySQLDSN_Default(t *testing.T) {
	got := portFromURL("user:pass@tcp(localhost:3306)/db", "3306")
	if got != "3306" {
		t.Errorf("expected 3306, got %s", got)
	}
}
