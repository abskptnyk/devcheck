package reporter

import (
	"encoding/json"
	"fmt"

	"github.com/vidya381/devcheck/internal/check"
)

type jsonResult struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Fix     string `json:"fix,omitempty"`
}

func RenderJSON(results []check.Result) {
	out := make([]jsonResult, len(results))
	for i, r := range results {
		status := "pass"
		switch r.Status {
		case check.StatusWarn:
			status = "warn"
		case check.StatusFail:
			status = "fail"
		case check.StatusSkipped:
			status = "skipped"
		}
		out[i] = jsonResult{
			Name:    r.Name,
			Status:  status,
			Message: r.Message,
			Fix:     r.Fix,
		}
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(b))
}
