package dotlink

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func exists(fn string) (bool, error) {
	if _, err := os.Stat(fn); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, fmt.Errorf("os.Stat: %w", err)
	}

	return true, nil
}

func Import(root, target string, ignores PathMatcher) (States, error) {
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

		exists, err := exists(target)
		if err != nil {
			return fmt.Errorf("exists: %w", err)
		}

		if !exists || ignores.MatchesPath(strings.TrimPrefix(path, root)) {
			return nil
		}

		s = append(s, State{
			From: path,
			To:   target,
		})

		return nil
	})

	return s, nil
}
