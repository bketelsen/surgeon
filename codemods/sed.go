package codemods

import (
	"errors"
	"fmt"
	"log/slog"
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
	slog.Info("Applying sed", "source", source, "target", target, "match", match, "args", args)

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

func (s Sed) Validate(_, _, _ string, args ...string) error {
	if len(args) != 2 {
		return errors.New("sed requires two arguments")
	}
	return nil
}

func (s Sed) Description() string {
	return "Replace strings in a file"
}

func (s Sed) Usage() string {
	return `Replace a strings in a file.
This codemod replaces strings in the matched file(s) with a string
specified in the arguments.

Args (2 required):
	1. search string
	2. replacement string

Example:
	upstream: https://github.com/community-scripts/ProxmoxVE
	modsdir: codemods
	codemods:
	- description: Header Updates
		mod: sed
		match: misc/*.func
		args:
		- https://github.com/community-scripts/ProxmoxVE/raw/main/ct/headers/
		- https://github.com/bketelsen/IncusScripts/raw/main/ct/headers/
	`
}

func sed(old, newthing, filePath string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	fi, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	fileString := string(fileData)
	fileString = strings.ReplaceAll(fileString, old, newthing)
	fileData = []byte(fileString)

	err = os.WriteFile(filePath, fileData, fi.Mode())
	if err != nil {
		return err
	}

	return nil
}
