package checks

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/vidya381/devcheck/internal/check"
)

type BinaryCheck struct {
	Binary string
}

func (c *BinaryCheck) Name() string {
	return fmt.Sprintf("%s installed", c.Binary)
}

func (c *BinaryCheck) Run(_ context.Context) check.Result {
	_, err := exec.LookPath(c.Binary)
	if err != nil {
		return check.Result{
			Name:    c.Name(),
			Status:  check.StatusFail,
			Message: fmt.Sprintf("%s not found on PATH", c.Binary),
			Fix:     fmt.Sprintf("Install %s and make sure it is on your PATH", c.Binary),
		}
	}
	return check.Result{
		Name:    c.Name(),
		Status:  check.StatusPass,
		Message: fmt.Sprintf("%s is installed", c.Binary),
	}
}
