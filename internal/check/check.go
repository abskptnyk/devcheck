package check

import "context"

type Status int

const (
	StatusPass Status = iota
	StatusWarn
	StatusFail
	StatusSkipped
)

type Result struct {
	Name    string
	Status  Status
	Message string
	Fix     string // shown when --fix flag is passed
}

type Check interface {
	Name() string
	Run(ctx context.Context) Result
}
