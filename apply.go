package dotlink

import "fmt"

func Apply(current States, root, target string, ignores PathMatcher) error {
	diffs, err := Plan(current, root, target, ignores)
	if err != nil {
		return fmt.Errorf("Plan: %w", err)
	}

	if err := diffs.Apply(); err != nil {
		return fmt.Errorf("diffs.Apply: %w", err)
	}

	return nil
}
