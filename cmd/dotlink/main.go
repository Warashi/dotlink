package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
	case "apply":
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
	defer f.Close()

	if err := states.SaveTo(f); err != nil {
		return fmt.Errorf("states.SaveTo: %w", err)
	}

	return nil
}
