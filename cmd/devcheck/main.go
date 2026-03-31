package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/vidya381/devcheck/internal/check"
	"github.com/vidya381/devcheck/internal/config"
	"github.com/vidya381/devcheck/internal/detector"
	"github.com/vidya381/devcheck/internal/reporter"
)

var version = "dev"

var (
	flagVerbose bool
	flagJSON    bool
	flagFix     bool
	flagCI      bool
)

func main() {
	root := &cobra.Command{
		Use:     "devcheck",
		Short:   "Check if your dev environment is ready to run this project",
		Version: version,
		RunE:    run,
	}

	root.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, "Show all checks including skipped")
	root.Flags().BoolVar(&flagJSON, "json", false, "Output results as JSON")
	root.Flags().BoolVar(&flagFix, "fix", false, "Show suggested fix for each failure")
	root.Flags().BoolVar(&flagCI, "ci", false, "Exit with code 1 on any failure (for CI pipelines)")

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	dir, _ := os.Getwd()

	cfg, err := config.Load(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}

	stack := detector.Detect(dir)
	checks := check.Build(stack)

	// Append extra binary checks from devcheck.yml require:
	for _, binary := range cfg.Require {
		checks = append(checks, &check.BinaryCheck{Binary: binary})
	}

	// Filter out checks the user asked to skip:
	if len(cfg.Skip) > 0 {
		skipSet := cfg.SkipSet()
		filtered := checks[:0]
		for _, c := range checks {
			if _, skip := skipSet[c.Name()]; !skip {
				filtered = append(filtered, c)
			}
		}
		checks = filtered
	}

	results := make([]check.Result, len(checks))
	var wg sync.WaitGroup

	for i, c := range checks {
		wg.Add(1)
		go func(i int, c check.Check) {
			defer wg.Done()
			results[i] = c.Run(context.Background())
		}(i, c)
	}
	wg.Wait()

	if flagJSON {
		reporter.RenderJSON(results)
	} else {
		reporter.Render(results, flagFix)
	}

	if flagCI {
		for _, r := range results {
			if r.Status == check.StatusFail {
				os.Exit(1)
			}
		}
	}

	return nil
}
