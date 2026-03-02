package check

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// GoVersionCheck reads the required Go version from go.mod and compares
// it against the installed version.
type GoVersionCheck struct {
	Dir string
}

func (c *GoVersionCheck) Name() string {
	return "Go version"
}

func (c *GoVersionCheck) Run(_ context.Context) Result {
	required, err := readGoModVersion(c.Dir + "/go.mod")
	if err != nil {
		return Result{Name: c.Name(), Status: StatusSkipped, Message: "could not read go.mod"}
	}

	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return Result{Name: c.Name(), Status: StatusFail, Message: "could not run go version"}
	}
	// output: "go version go1.22.3 darwin/arm64"
	parts := strings.Fields(string(out))
	if len(parts) < 3 {
		return Result{Name: c.Name(), Status: StatusFail, Message: "unexpected go version output"}
	}
	installed := strings.TrimPrefix(parts[2], "go")

	if versionLess(installed, required) {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: fmt.Sprintf("need Go %s, got %s", required, installed),
			Fix:     fmt.Sprintf("upgrade Go to %s or newer", required),
		}
	}
	return Result{
		Name:    c.Name(),
		Status:  StatusPass,
		Message: fmt.Sprintf("Go %s installed (need %s)", installed, required),
	}
}

// NodeVersionCheck reads the required Node version from .nvmrc or
// engines.node in package.json and compares against the installed version.
type NodeVersionCheck struct {
	Dir string
}

func (c *NodeVersionCheck) Name() string {
	return "Node version"
}

func (c *NodeVersionCheck) Run(_ context.Context) Result {
	required, err := readNodeRequired(c.Dir)
	if err != nil {
		return Result{Name: c.Name(), Status: StatusSkipped, Message: "no Node version requirement found"}
	}

	out, err := exec.Command("node", "--version").Output()
	if err != nil {
		return Result{Name: c.Name(), Status: StatusFail, Message: "could not run node --version"}
	}
	installed := strings.TrimPrefix(strings.TrimSpace(string(out)), "v")

	if versionLess(installed, required) {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: fmt.Sprintf("need Node %s, got %s", required, installed),
			Fix:     fmt.Sprintf("upgrade Node to %s or newer", required),
		}
	}
	return Result{
		Name:    c.Name(),
		Status:  StatusPass,
		Message: fmt.Sprintf("Node %s installed (need %s)", installed, required),
	}
}

func readGoModVersion(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "go ") {
			return strings.TrimPrefix(line, "go "), nil
		}
	}
	return "", fmt.Errorf("go directive not found in go.mod")
}

func readNodeRequired(dir string) (string, error) {
	// .nvmrc takes priority
	if data, err := os.ReadFile(dir + "/.nvmrc"); err == nil {
		v := strings.TrimSpace(string(data))
		return strings.TrimPrefix(v, "v"), nil
	}

	// fall back to engines.node in package.json
	data, err := os.ReadFile(dir + "/package.json")
	if err != nil {
		return "", err
	}
	var pkg struct {
		Engines struct {
			Node string `json:"node"`
		} `json:"engines"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return "", err
	}
	if pkg.Engines.Node == "" {
		return "", fmt.Errorf("no engines.node in package.json")
	}
	// strip range operators like >=, ^, ~
	return strings.TrimLeft(pkg.Engines.Node, ">=^~"), nil
}

// versionLess returns true if a < b.
func versionLess(a, b string) bool {
	av := parseVersion(a)
	bv := parseVersion(b)
	for i := 0; i < len(bv); i++ {
		if i >= len(av) {
			return true
		}
		if av[i] < bv[i] {
			return true
		}
		if av[i] > bv[i] {
			return false
		}
	}
	return false
}

func parseVersion(v string) []int {
	var result []int
	for _, p := range strings.Split(v, ".") {
		n, err := strconv.Atoi(p)
		if err != nil {
			break
		}
		result = append(result, n)
	}
	return result
}
