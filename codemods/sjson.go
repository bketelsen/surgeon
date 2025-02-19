package codemods

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/tidwall/sjson"
)

func init() {
	Mods["sjson"] = SJSON{}
}

type SJSON struct{}

// assert that Inject implements CodeMod
var _ CodeMod = SJSON{}

func (s SJSON) Apply(source, target, match string, args ...string) error {

	slog.Info("Applying sjson", "source", source, "target", target, "match", match, "args", args)

	sourceMatches := filepath.Join(source, match)
	matches, err := filepath.Glob(sourceMatches)
	if err != nil {
		return fmt.Errorf("globbing source: %w", err)
	}
	for _, m := range matches {
		action := args[0]
		key := args[1]
		var value string
		if len(args) == 3 {
			value = args[2]
		}

		output, err := apply(action, key, value, m)
		if err != nil {
			return fmt.Errorf("modifying json: %w", err)
		}

		err = os.WriteFile(m, []byte(output), 0644)
		if err != nil {
			return fmt.Errorf("writing file: %w", err)
		}

	}

	return nil
}
func (s SJSON) Validate(source, target, match string, args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("sjson requires at least two arguments")
	}
	return nil
}

func (s SJSON) Description() string {
	return "Modify a JSON file in-place"
}

func (s SJSON) Usage() string {
	return `sjson modifies a JSON file in-place.
This codemod modifies the matched file(s) by injecting specified content.

Args (3 required for set, 2 required for del):
	1. Action (set, del)
	2. Key path
	3. Value (required for set)

Example:
	upstream: https://github.com/community-scripts/ProxmoxVE
	modsdir: codemods
	codemods:
	- description: change OS Key to debian
		mod: sjson
		match: json/debian-vm.json
		args:
		- set
		- install_methods.1.resources.os
		- debian
	`
}

func apply(action, key, value, filePath string) (string, error) {

	bb, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	content := string(bb)
	var output string
	switch action {
	case "set":
		slog.Debug("Setting", "key", key, "to", value)
		output, err = sjson.Set(content, key, value)
		if err != nil {
			return "", err
		}
		//fmt.Println("Output", output)

	case "del":
		// add newline
		slog.Debug("Deleting", "key", key)
		output, err = sjson.Delete(content, key)
		if err != nil {
			return "", err
		}

	default:
		return "", fmt.Errorf("unknown action %q", action)

	}

	return output, nil
}
