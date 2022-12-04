package dotlink

import (
	"fmt"
	"os"
)

func exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("os.Stat: %w", err)
	}
	return true, nil
}

func Import(root, target string, ignores PathMatcher) (States, error) {
	s, err := NewStates(root, target, ignores)
	if err != nil {
		return nil, fmt.Errorf("NewStates: %w", err)
	}

	for i := range s {
		e, err := exists(s[i].To)
		if err != nil {
			return nil, fmt.Errorf("exists: %w", err)
		}

		// 存在していない == 無視されている == Ignored
		s[i].Ignored = !e
	}

	return s, nil
}
