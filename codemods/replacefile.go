package codemods

import (
	"fmt"
	"os"
	"path/filepath"
)

func init() {
	Mods["replacefile"] = ReplaceFile{}
}

type ReplaceFile struct{}

// assert that Sed implements CodeMod
var _ CodeMod = ReplaceFile{}

func (s ReplaceFile) Apply(source, target, match string, args ...string) error {
	fmt.Println("Source", source)
	fmt.Println("Target", target)
	fmt.Println("Match", match)
	fmt.Println("Args", args)
	sourceMatches := filepath.Join(source, match)
	matches, err := filepath.Glob(sourceMatches)
	if err != nil {
		return fmt.Errorf("globbing source: %w", err)
	}
	for _, m := range matches {
		replacementPath := filepath.Join(target, args[0])
		err = replace(replacementPath, m)
		if err != nil {
			return fmt.Errorf("applying replacement: %w", err)
		}
	}

	return nil
}
func (s ReplaceFile) Validate(source, target, match string, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("replacefile requires two arguments")
	}
	return nil
}

func replace(new, filePath string) error {
	fmt.Println("Replacing", filePath, "with", new)
	bb, err := os.ReadFile(new)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, bb, 0644)
	if err != nil {
		return err
	}
	return nil
}
