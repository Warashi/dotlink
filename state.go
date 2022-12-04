package dotlink

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/google/go-cmp/cmp"
)

type (
	State struct {
		From, To string
		Ignored  bool
	}
	States     []State
	StateDiff  [2]State
	StateDiffs []StateDiff
)

func (s State) Less(o State) bool {
	if s.From != o.From {
		return s.From < o.From
	}

	if s.To != o.To {
		return s.To < o.To
	}

	if s.Ignored != o.Ignored {
		// ignored が異なる場合には true が優先
		// l.Ignored == true なら l が優先
		// l.Ignored == false なら r.Ignored == true なので r が優先
		return s.Ignored
	}

	// 全く同じならとりあえず true を返す
	return true

}

func (s States) Len() int           { return len(s) }
func (s States) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s States) Less(i, j int) bool { return s[i].Less(s[j]) }

func (s States) SaveTo(w io.Writer) error {
	if err := json.NewEncoder(w).Encode(s); err != nil {
		return fmt.Errorf("json.NewEncoder.Encode: %w", err)
	}
	return nil
}

func LoadState(r io.Reader) (States, error) {
	var s States
	if err := json.NewDecoder(r).Decode(&s); err != nil {
		return nil, fmt.Errorf("json.NewDecoder.Decode: %w", err)
	}
	return s, nil
}

func (s States) Clone() States {
	c := make(States, len(s))
	copy(c, s)
	return c
}

func (s States) Diff(o States) (StateDiffs, error) {
	result := make(StateDiffs, 0, len(s)+len(o))

	sset := make(map[string]State, len(s))
	for _, s := range s {
		sset[s.From] = s
	}
	for _, o := range o {
		s := sset[o.From]
		if s != o {
			result = append(result, StateDiff{s, o})
		}
	}

	oset := make(map[string]State, len(s))
	for _, o := range o {
		oset[o.From] = o
	}
	for _, s := range s {
		o := oset[s.From]
		if s != o {
			result = append(result, StateDiff{s, o})
		}
	}

	return result, nil
}

func (d StateDiff) String() string {
	return cmp.Diff(d[0], d[1])
}

func (d StateDiff) Less(o StateDiff) bool {
	if d[0] != o[0] {
		return d[0].Less(o[0])
	}
	if d[1] != o[1] {
		return d[1].Less(o[1])
	}
	return true
}

func (d StateDiff) Apply() error {
	o, n := d[0], d[1]
	if o.Ignored && !n.Ignored {
		if err := os.Symlink(n.From, n.To); err != nil {
			return fmt.Errorf("os.Symlink: %w", err)
		}
		return nil
	}
	if !o.Ignored && n.Ignored {
		if err := os.Remove(o.To); err != nil {
			return fmt.Errorf("os.Remove", err)
		}
		return nil
	}

	if o.To != n.To {
		if err := os.Remove(o.To); err != nil {
			return fmt.Errorf("os.Remove", err)
		}
		if err := os.Symlink(n.From, n.To); err != nil {
			return fmt.Errorf("os.Symlink: %w", err)
		}

		return nil
	}

	return nil
}

func (d StateDiffs) Len() int           { return len(d) }
func (d StateDiffs) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d StateDiffs) Less(i, j int) bool { return d[i].Less(d[j]) }

func (d StateDiffs) Clone() StateDiffs {
	c := make(StateDiffs, len(d))
	copy(c, d)
	return c
}

func (d StateDiffs) Apply() error {
	d = d.Clone()
	sort.Sort(d)
	for _, d := range d {
		if err := d.Apply(); err != nil {
			return fmt.Errorf("d.Apply: %w", err)
		}
	}

	return nil
}
