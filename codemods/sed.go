package codemods

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	Mods["sed"] = Sed{}
}

type Sed struct{}

// assert that Sed implements CodeMod
var _ CodeMod = Sed{}

func (s Sed) Apply(source, target, match string, args ...string) error {
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
		err = sed(args[0], args[1], m)
		if err != nil {
			return fmt.Errorf("applying sed: %w", err)
		}
	}

	return nil
}
func (s Sed) Validate(source, target, match string, args ...string) error {
	if len(args) != 2 {
		return fmt.Errorf("sed requires two arguments")
	}
	return nil
}

func sed(old, new, filePath string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	fi, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	fileString := string(fileData)
	fileString = strings.ReplaceAll(fileString, old, new)
	fileData = []byte(fileString)

	err = os.WriteFile(filePath, fileData, fi.Mode())
	if err != nil {
		return err
	}

	return nil
}
