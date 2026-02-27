package check

import "github.com/vidya381/devcheck/internal/detector"

func Build(stack detector.DetectedStack) []Check {
	var checks []Check

	if stack.Go {
		// add Go checks
	}
	if stack.Node {
		// add Node checks
	}
	if stack.Python {
		// add Python checks
	}
	if stack.Java {
		// add Java checks
	}
	if stack.Docker {
		// add Docker checks
	}
	if stack.Postgres {
		// add Postgres checks
	}
	if stack.Redis {
		// add Redis checks
	}

	// always run env check if .env.example exists
	// checks = append(checks, &EnvCheck{})

	return checks
}
