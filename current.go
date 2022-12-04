package dotlink

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func NewStates(root, target string, ignores PathMatcher) (States, error) {
	var s States
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		target := strings.Replace(path, root, target, 1)
		s = append(s, State{
			From:    path,
			To:      target,
			Ignored: ignores.MatchesPath(path),
		})

		return nil
	})

	return s, nil
}
