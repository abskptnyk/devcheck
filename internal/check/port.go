package check

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

type PortCheck struct {
	Service string
	Port    string
	dialer  func(address string) error
}

func (c *PortCheck) Name() string {
	return fmt.Sprintf("%s port reachable", c.Service)
}

func (c *PortCheck) Run(_ context.Context) Result {
	dial := c.dialer
	if dial == nil {
		dial = func(address string) error {
			conn, err := net.DialTimeout("tcp", address, 2*time.Second)
			if err != nil {
				return err
			}
			conn.Close()
			return nil
		}
	}

	address := "localhost:" + c.Port
	if err := dial(address); err != nil {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: fmt.Sprintf("nothing listening on port %s", c.Port),
			Fix:     fmt.Sprintf("start %s and make sure it is running on port %s", c.Service, c.Port),
		}
	}
	return Result{
		Name:    c.Name(),
		Status:  StatusPass,
		Message: fmt.Sprintf("%s is listening on port %s", c.Service, c.Port),
	}
}

// portFromURL extracts the port from a service URL, falling back to defaultPort.
// Handles standard URLs (postgres://, redis://, mongodb://) and MySQL DSNs (user:pass@tcp(host:port)/db).
func portFromURL(rawURL, defaultPort string) string {
	if rawURL == "" {
		return defaultPort
	}

	// MySQL DSN format: user:pass@tcp(host:port)/db
	if idx := strings.Index(rawURL, "tcp("); idx != -1 {
		rest := rawURL[idx+4:]
		if end := strings.Index(rest, ")"); end != -1 {
			_, port, err := net.SplitHostPort(rest[:end])
			if err == nil && port != "" {
				return port
			}
		}
	}

	// Standard URL format
	u, err := url.Parse(rawURL)
	if err == nil && u.Port() != "" {
		return u.Port()
	}

	return defaultPort
}
