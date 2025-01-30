package codemods

import (
	"fmt"
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
		replacement := filepath.Join(target, args[1])
		err = replaceFunction(args[0], replacement, m)
		if err != nil {
			return fmt.Errorf("applying bash function replacer: %w", err)
		}
	}

	return nil
}
func (s BashFunc) Validate(source, target, match string, args ...string) error {
	if len(args) != 2 {
		return fmt.Errorf("bashfunc requires two arguments")
	}
	return nil
}

func replaceFunction(name, replacementPath, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	f, err := syntax.NewParser().Parse(file, "")
	if err != nil {
		return err
	}
	file.Close()
	var startpos uint
	var endpos uint
	var found bool
	for _, stmt := range f.Stmts {
		decl, ok := stmt.Cmd.(*syntax.FuncDecl)
		if ok {
			if decl.Name.Value == name {
				startpos = decl.Pos().Line()
				endpos = decl.End().Line()
				found = true
				break
			}
		}
	}
	if found {
		// // read the replacement file
		replacement, err := os.Open(replacementPath)
		defer replacement.Close()
		if err != nil {
			fmt.Println("Error: ", err)
			return err
		}
		replace, err := syntax.NewParser().Parse(replacement, "")
		if err != nil {
			return err
		}

		bb, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error: ", err)
			return err
		}

		// slice the file into 3 parts
		// 1. before the function
		// 2. the function
		// 3. after the function
		contents := string(bb)
		lines := strings.Split(contents, "\n")
		var before []string
		var after []string
		var inside []string

		for i, line := range lines {
			if i < int(startpos-1) {
				before = append(before, line)
			} else if i >= int(startpos-1) && i < int(endpos) {
				inside = append(inside, line)
			} else {
				after = append(after, line)
			}
		}

		// write the before part
		nf, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error: ", err)
			return err
		}
		defer nf.Close()
		for _, line := range before {
			nf.WriteString(line + "\n")
		}
		// write the replacement
		syntax.NewPrinter().Print(nf, replace)
		nf.WriteString("\n")
		// write the after part
		for _, line := range after {
			nf.WriteString(line + "\n")
		}
		nf.Close()

	}
	return nil
}
