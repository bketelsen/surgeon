// Package codemods provides a set of code modifications (codemods) for
// modifying files. It allows you to define a set of code modifications
// and apply them to a directory.
package codemods

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func init() {
	Mods["inject"] = Inject{}
}

type Inject struct{}

// assert that Inject implements CodeMod
var _ CodeMod = Inject{}

func (s Inject) Apply(source, target, match string, args ...string) error {
	slog.Info("Applying code injector", "source", source, "target", target, "match", match, "args", args)

	sourceMatches := filepath.Join(source, match)
	matches, err := filepath.Glob(sourceMatches)
	if err != nil {
		return fmt.Errorf("globbing source: %w", err)
	}
	for _, m := range matches {
		where := args[0]
		contents := args[1]
		err = inject(where, contents, m)
		if err != nil {
			return fmt.Errorf("injecting content: %w", err)
		}
	}

	return nil
}

func (s Inject) Validate(_, _, _ string, args ...string) error {
	if len(args) != 2 {
		return errors.New("inject requires two arguments")
	}
	return nil
}

func (s Inject) Description() string {
	return "Inject contents into a file"
}

func (s Inject) Usage() string {
	return `Inject contents into a file.
This codemod modifies the matched file(s) by injecting specified content.

Args (2 required):
	1. Injection point in the file. Valid: "start", "end", <line number>
	2. The content to inject

Example:
	upstream: https://github.com/community-scripts/ProxmoxVE
	modsdir: codemods
	codemods:
	- description: Inject Modification notice
		mod: inject
		match: install/*.sh
		args:
		- end
		- # Modified by surgeon
	`
}

func inject(where, contents, filePath string) error {
	slog.Debug("Injecting", "contents", contents, "at", where, "in", filePath)
	bb, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	switch where {
	case "start":
		// add newline
		contents += "\n"
		bb = append([]byte(contents), bb...)
	case "end":
		// add newline
		bb = append(bb, []byte("\n")...)
		bb = append(bb, []byte(contents)...)
	default:
		line, err := strconv.Atoi(where)
		if err != nil {
			return fmt.Errorf("invalid line number: %w", err)
		}
		strcontent := string(bb)
		lines := strings.Split(strcontent, "\n")
		if line > len(lines) {
			return errors.New("line number out of range")
		}
		var buf bytes.Buffer
		for i, l := range lines {
			buf.WriteString(l)
			if i == line {
				buf.WriteString(contents)
			}
			buf.WriteString("\n")
		}
		bb = buf.Bytes()
	}

	err = os.WriteFile(filePath, bb, 0o644)
	if err != nil {
		return err
	}
	return nil
}
