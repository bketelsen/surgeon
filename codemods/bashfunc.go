package codemods

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"mvdan.cc/sh/syntax"
)

func init() {
	Mods["bashfunc"] = BashFunc{}
}

type BashFunc struct{}

// assert that Sed implements CodeMod
var _ CodeMod = BashFunc{}

func (s BashFunc) Apply(source, target, match string, args ...string) error {
	slog.Info("Applying bash function replacer", "source", source, "target", target, "match", match, "args", args)
	sourceMatches := filepath.Join(source, match)
	matches, err := filepath.Glob(sourceMatches)
	if err != nil {
		return fmt.Errorf("globbing source: %w", err)
	}
	for _, m := range matches {
		replacement := filepath.Join(target, args[1])
		err = replaceFunctionInFile(args[0], replacement, m)
		if err != nil {
			return fmt.Errorf("applying bash function replacer: %w", err)
		}
	}

	return nil
}

func (s BashFunc) Validate(_, _, _ string, args ...string) error {
	if len(args) != 2 {
		return errors.New("bashfunc requires two arguments")
	}
	return nil
}

func (s BashFunc) Description() string {
	return "Replace a bash function with another"
}

func (s BashFunc) Usage() string {
	return `Replace a bash function with another.
This codemod searches for a bash function in the matched file(s)
and replaces it with another function.

Args (2 required):
	1. The name of the function to replace
	2. The path to the file (in your fork) containing the replacement function

Example:
	upstream: https://github.com/community-scripts/ProxmoxVE
	modsdir: codemods
	codemods:
	- description: PVE Check Function
		mod: bashfunc
		match: misc/build.func
		args:
		- pve_check
		- codemods/pve_check.sh
	`
}

func replaceFunction(name string, replacementContent, fileContent []byte) ([]byte, error) {
	f, err := syntax.NewParser().Parse(bytes.NewReader(fileContent), "")
	if err != nil {
		return nil, err
	}

	var startpos uint
	var endpos uint
	var found bool
	for _, stmt := range f.Stmts {
		decl, ok := stmt.Cmd.(*syntax.FuncDecl)
		if ok && decl.Name.Value == name {
			startpos = decl.Pos().Line()
			endpos = decl.End().Line()
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("function %q not found", name)
	}

	// Parse the replacement content
	replace, err := syntax.NewParser().Parse(bytes.NewReader(replacementContent), "")
	if err != nil {
		return nil, err
	}

	// Slice the file content into 3 parts: before, function, and after
	lines := strings.Split(string(fileContent), "\n")
	var before []string
	var after []string

	for i, line := range lines {
		if i < int(startpos-1) {
			before = append(before, line)
		} else if i >= int(startpos-1) && i < int(endpos) {
			continue
		} else {
			after = append(after, line)
		}
	}

	// Combine the parts with the replacement
	var result bytes.Buffer
	for _, line := range before {
		result.WriteString(line + "\n")
		result.WriteString(line)
	}
	syntax.NewPrinter().Print(&result, replace)
	// result.WriteString("\n")
	for _, line := range after {
		result.WriteString(line + "\n")
	}
	// Remove the last newline character if it exists
	if result.Len() > 0 && result.Bytes()[result.Len()-1] == '\n' {
		result.Truncate(result.Len() - 1)
	}
	return result.Bytes(), nil
}

func replaceFunctionInFile(name, replacementPath, filePath string) error {
	// Read the original file content
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Read the replacement content
	replacementContent, err := os.ReadFile(replacementPath)
	if err != nil {
		return err
	}

	// Perform the replacement
	modifiedContent, err := replaceFunction(name, replacementContent, fileContent)
	if err != nil {
		return err
	}

	// Write the modified content back to the file
	err = os.WriteFile(filePath, modifiedContent, 0o644)
	if err != nil {
		return err
	}

	return nil
}
