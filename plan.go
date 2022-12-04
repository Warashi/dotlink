package dotlink

import "fmt"

func Plan(current States, root, target string, ignores PathMatcher) (StateDiffs, error) {
	snew, err := NewStates(root, target, ignores)
	if err != nil {
		return nil, fmt.Errorf("NewStates: %w", err)
	}

	diff, err := current.Diff(snew)
	if err != nil {
		return nil, fmt.Errorf("current.Diff: %w", err)
	}

	return diff, nil
}
