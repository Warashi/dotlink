package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Warashi/dotlink"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("os.UserHomeDir: %v", err)
	}

	var root, target, ignores, state string
	flag.StringVar(&root, "root", ".", "root for files (default: current directory)")
	flag.StringVar(&target, "target", home, "target root for files (default: home directory)")
	flag.StringVar(&ignores, "ignores", "./.dtlignore", "ignore files (default: ./.dtlignore)")
	flag.StringVar(&state, "state", "./state.json", "state file (default: ./state.json)")

	flag.Parse()

	os.Exit(_main(flag.Arg(0), root, target, ignores, state))
}

func _main(cmd, root, target, ignores, state string) int {
	switch cmd {
	case "", "help":
		flag.Usage()
		return 1
	case "import":
		if err := Import(root, target, ignores, state); err != nil {
			fmt.Printf("dotlink: import: %v", err)
			return 1
		}
	case "plan":
		if err := Plan(root, target, ignores, state); err != nil {
			fmt.Printf("dotlink: plan: %v", err)
			return 1
		}
	case "apply":
		if err := Apply(root, target, ignores, state); err != nil {
			fmt.Printf("dotlink: apply: %v", err)
			return 1
		}
	case "destroy":
		if err := Destroy(root, target, ignores, state); err != nil {
			fmt.Printf("dotlink: destroy: %v", err)
			return 1
		}
	default:
		fmt.Printf("dotlink: unknown command: %s\n", cmd)
	}

	return 0
}

func Import(root, target, ignores, state string) error {
	ig, err := dotlink.ParseIgnores(ignores)
	if err != nil {
		return fmt.Errorf("dotlink.ParseIgnores: %w", err)
	}

	states, err := dotlink.Import(root, target, ig)
	if err != nil {
		return fmt.Errorf("dotlink.Import: %w", err)
	}

	f, err := os.Create(state)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}

	if err := states.SaveTo(f); err != nil {
		return fmt.Errorf("states.SaveTo: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("f.Close: %w", err)
	}

	return nil
}

func Plan(root, target, ignores, state string) error {
	ig, err := dotlink.ParseIgnores(ignores)
	if err != nil {
		return fmt.Errorf("dotlink.ParseIgnores: %w", err)
	}

	f, err := os.Open(state)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}

	states, err := dotlink.LoadState(f)
	if err != nil {
		return fmt.Errorf("dotlink.LoadState: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("f.Close: %w", err)
	}

	diffs, err := dotlink.Plan(states, root, target, ig)
	if err != nil {
		return fmt.Errorf("dotlink.Import: %w", err)
	}

	fmt.Println(diffs)

	return nil
}

func Apply(root, target, ignores, state string) error {
	ig, err := dotlink.ParseIgnores(ignores)
	if err != nil {
		return fmt.Errorf("dotlink.ParseIgnores: %w", err)
	}

	f, err := os.Open(state)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}

	states, err := dotlink.LoadState(f)
	if err != nil {
		return fmt.Errorf("dotlink.LoadState: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("f.Close: %w", err)
	}

	diffs, err := dotlink.Plan(states, root, target, ig)
	if err != nil {
		return fmt.Errorf("dotlink.Import: %w", err)
	}

	fmt.Println(diffs)

	fmt.Println(strings.Repeat("#", 20))
	fmt.Println("# applying")
	fmt.Println(strings.Repeat("#", 20))

	if err := diffs.Apply(); err != nil {
		return fmt.Errorf("diffs.Apply: %w", err)
	}

	if err := Import(root, target, ignores, state); err != nil {
		return fmt.Errorf("Import: %w", err)
	}

	return nil
}

func Destroy(root, target, ignores, state string) error {
	f, err := os.Open(state)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}

	states, err := dotlink.LoadState(f)
	if err != nil {
		return fmt.Errorf("dotlink.LoadState: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("f.Close: %w", err)
	}

	if err := dotlink.Destroy(states); err != nil {
		return fmt.Errorf("dotlink.Destroy: %w", err)
	}

	if err := Import(root, target, ignores, state); err != nil {
		return fmt.Errorf("Import: %w", err)
	}

	return nil
}
