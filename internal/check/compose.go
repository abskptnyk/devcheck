package check

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type ComposeCheck struct {
	runner func() ([]byte, error)
}

func (c *ComposeCheck) Name() string {
	return "Docker Compose services running"
}

type composeService struct {
	Name  string `json:"Name"`
	State string `json:"State"`
}

func (c *ComposeCheck) Run(_ context.Context) Result {
	run := c.runner
	if run == nil {
		run = func() ([]byte, error) {
			return exec.Command("docker", "compose", "ps", "--format", "json").Output()
		}
	}

	out, err := run()
	if err != nil {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: "could not run docker compose ps",
			Fix:     "make sure Docker is running and you are in the project directory",
		}
	}

	// docker compose ps --format json outputs one JSON object per line
	var stopped []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var svc composeService
		if err := json.Unmarshal([]byte(line), &svc); err != nil {
			continue
		}
		if svc.State != "running" {
			stopped = append(stopped, svc.Name)
		}
	}

	if len(stopped) > 0 {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: fmt.Sprintf("services not running: %s", strings.Join(stopped, ", ")),
			Fix:     "run docker compose up -d to start them",
		}
	}

	return Result{
		Name:    c.Name(),
		Status:  StatusPass,
		Message: "all services are running",
	}
}

type ComposeImageCheck struct {
	runner func() ([]byte, error)
}

func (c *ComposeImageCheck) Name() string {
	return "Docker Compose images pulled"
}

type composeImage struct {
	ContainerName string `json:"ContainerName"`
	Repository    string `json:"Repository"`
	ID            string `json:"ID"`
}

func (c *ComposeImageCheck) Run(_ context.Context) Result {
	run := c.runner
	if run == nil {
		run = func() ([]byte, error) {
			return exec.Command("docker", "compose", "images", "--format", "json").Output()
		}
	}

	out, err := run()
	if err != nil {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: "could not run docker compose images",
			Fix:     "make sure Docker is running and you are in the project directory",
		}
	}

	images, parseErr := parseComposeImages(out)
	if parseErr != nil {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: "could not parse docker compose images output",
		}
	}

	var missing []string
	for _, img := range images {
		if img.Repository == "" || img.Repository == "<none>" || img.ID == "" {
			missing = append(missing, img.ContainerName)
		}
	}

	if len(missing) > 0 {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: fmt.Sprintf("images not pulled for: %s", strings.Join(missing, ", ")),
			Fix:     "run docker compose pull to pull all images",
		}
	}

	return Result{
		Name:    c.Name(),
		Status:  StatusPass,
		Message: "all service images are pulled",
	}
}

func parseComposeImages(data []byte) ([]composeImage, error) {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil, nil
	}
	// newer Docker versions output a JSON array; older versions output JSONL
	if data[0] == '[' {
		var images []composeImage
		if err := json.Unmarshal(data, &images); err != nil {
			return nil, err
		}
		return images, nil
	}
	var images []composeImage
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var img composeImage
		if err := json.Unmarshal([]byte(line), &img); err != nil {
			continue
		}
		images = append(images, img)
	}
	return images, nil
}
