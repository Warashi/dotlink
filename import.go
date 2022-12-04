package dotlink

import (
	"fmt"
)

func Import(root, target string, ignores PathMatcher) (States, error) {
	s, err := NewStates(root, target, ignores)
	if err != nil {
		return nil, fmt.Errorf("NewStates: %w", err)
	}

	return s, nil
}
