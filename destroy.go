package dotlink

import "fmt"

func Destroy(current States) error {
	for _, s := range current {
		if err := s.Destroy(); err != nil {
			return fmt.Errorf("s.Destroy: %w", err)
		}
	}

	return nil
}
