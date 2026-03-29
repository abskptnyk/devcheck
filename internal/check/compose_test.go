package check

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestComposeCheck_AllRunning(t *testing.T) {
	c := &ComposeCheck{runner: func() ([]byte, error) {
		return []byte(`{"Name":"app-db-1","State":"running"}` + "\n" +
			`{"Name":"app-web-1","State":"running"}` + "\n"), nil
	}}
	result := c.Run(context.Background())
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %v: %s", result.Status, result.Message)
	}
}

func TestComposeCheck_SomeStopped(t *testing.T) {
	c := &ComposeCheck{runner: func() ([]byte, error) {
		return []byte(`{"Name":"app-db-1","State":"running"}` + "\n" +
			`{"Name":"app-web-1","State":"exited"}` + "\n"), nil
	}}
	result := c.Run(context.Background())
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %v", result.Status)
	}
}

func TestComposeCheck_CommandFails(t *testing.T) {
	c := &ComposeCheck{runner: func() ([]byte, error) {
		return nil, errors.New("docker not found")
	}}
	result := c.Run(context.Background())
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %v", result.Status)
	}
}

func TestComposeCheck_NoServices(t *testing.T) {
	c := &ComposeCheck{runner: func() ([]byte, error) {
		return []byte(""), nil
	}}
	result := c.Run(context.Background())
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %v: %s", result.Status, result.Message)
	}
}

func TestComposeImageCheck_AllPulled(t *testing.T) {
	c := &ComposeImageCheck{runner: func() ([]byte, error) {
		return []byte(`[{"ContainerName":"app_web_1","Repository":"nginx","ID":"sha256:abc"}]`), nil
	}}
	result := c.Run(context.Background())
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %v: %s", result.Status, result.Message)
	}
}

func TestComposeImageCheck_MissingImage(t *testing.T) {
	c := &ComposeImageCheck{runner: func() ([]byte, error) {
		return []byte(`[{"ContainerName":"app_web_1","Repository":"","ID":""}]`), nil
	}}
	result := c.Run(context.Background())
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %v", result.Status)
	}
	if !strings.Contains(result.Message, "app_web_1") {
		t.Errorf("expected message to mention service name, got: %s", result.Message)
	}
}

func TestComposeImageCheck_NoneRepository(t *testing.T) {
	c := &ComposeImageCheck{runner: func() ([]byte, error) {
		return []byte(`[{"ContainerName":"app_db_1","Repository":"<none>","ID":""}]`), nil
	}}
	result := c.Run(context.Background())
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %v", result.Status)
	}
}

func TestComposeImageCheck_JSONLFormat(t *testing.T) {
	c := &ComposeImageCheck{runner: func() ([]byte, error) {
		return []byte(`{"ContainerName":"app_web_1","Repository":"nginx","ID":"sha256:abc"}` + "\n" +
			`{"ContainerName":"app_db_1","Repository":"postgres","ID":"sha256:def"}` + "\n"), nil
	}}
	result := c.Run(context.Background())
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %v: %s", result.Status, result.Message)
	}
}

func TestComposeImageCheck_CommandFails(t *testing.T) {
	c := &ComposeImageCheck{runner: func() ([]byte, error) {
		return nil, errors.New("docker not found")
	}}
	result := c.Run(context.Background())
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %v", result.Status)
	}
}

func TestComposeImageCheck_NoOutput(t *testing.T) {
	c := &ComposeImageCheck{runner: func() ([]byte, error) {
		return []byte(""), nil
	}}
	result := c.Run(context.Background())
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %v: %s", result.Status, result.Message)
	}
}
