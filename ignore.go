package dotlink

import (
	"fmt"

	ignore "github.com/sabhiram/go-gitignore"
)

type PathMatcher interface {
	MatchesPath(f string) bool
}

func ParseIgnores(fn string) (*ignore.GitIgnore, error) {
	i, err := ignore.CompileIgnoreFile(fn)
	if err != nil {
		return nil, fmt.Errorf("ignore.CompileIgnoreFile: %w", err)
	}
	return i, nil
}
