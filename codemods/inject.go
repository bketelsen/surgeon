package codemods

import (
	"bytes"
	"fmt"
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
		where := args[0]
		contents := args[1]
		err = inject(where, contents, m)
		if err != nil {
			return fmt.Errorf("injecting content: %w", err)
		}
	}

	return nil
}
func (s Inject) Validate(source, target, match string, args ...string) error {
	if len(args) != 2 {
		return fmt.Errorf("inject requires two arguments")
	}
	return nil
}

func (s Inject) Description() string {
	return "Inject contents into a file"
}

func (s Inject) Usage() string {
	return `Inject contents into a file.
This codemod modifies the matched file(s) by injecting specified content.

Args (3 required):
	1. Injection point in the file. Valid: "start", "end", line number
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
	fmt.Println("In", filePath, "with", contents, "at", where)
	bb, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	switch where {
	case "start":
		// add newline
		contents = contents + "\n"
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
			return fmt.Errorf("line number out of range")
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

	err = os.WriteFile(filePath, bb, 0644)
	if err != nil {
		return err
	}
	return nil
}
