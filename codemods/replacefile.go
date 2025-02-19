package codemods

import (
	"fmt"
	"log/slog"
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

	slog.Info("Applying replacefile", "source", source, "target", target, "match", match, "args", args)

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

func (s ReplaceFile) Description() string {
	return "Replace a file with another"
}

func (s ReplaceFile) Usage() string {
	return `Replace a file with another.
This codemod replaces the matched file(s) with a file from your fork.

Args (1 required):
	1. The path to the file (in your fork) to replace the matched file(s)

Example:
	upstream: https://github.com/community-scripts/ProxmoxVE
	modsdir: codemods
	codemods:
	- description: Replace create_lxc
		mod: replacefile
		match: ct/create_lxc.sh
		args:
		- codemods/create_lxc.sh
	`
}

func replace(new, filePath string) error {
	slog.Debug("Replacing", "file", filePath, "with", new)
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
