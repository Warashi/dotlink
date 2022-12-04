package dotlink

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func NewStates(root, target string, ignores PathMatcher) (States, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("filepath.Abs: %w", err)
	}

	var s States
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		target := strings.Replace(path, root, target, 1)
		if !ignores.MatchesPath(strings.TrimPrefix(path, root)) {
			s = append(s, State{
				From: path,
				To:   target,
			})
		}

		return nil
	})

	return s, nil
}
