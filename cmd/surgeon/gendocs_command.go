/*
Copyright © 2025 Brian Ketelsen <bketelsen@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bketelsen/surgeon/codemods"
	"github.com/bketelsen/toolbox/cobra"
	"github.com/bketelsen/toolbox/cobra/doc"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/viper"
)

// NewGendocsCommand creates a new command to generate documentation for the project
func NewGendocsCommand(config *viper.Viper) *cobra.Command {
	// gendocsCmd represents the gendocs command
	gendocsCmd := &cobra.Command{
		Use:    "gendocs",
		Hidden: true,
		Short:  "Generates documentation for the project",
		Long: `Generates documentation for the command using the cobra doc generator.
The documentation is generated in the directory specified by the --output flag and
is in markdown format.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			lipgloss.DefaultRenderer().SetColorProfile(termenv.Ascii)

			o := config.GetString("docs.output")
			cmd.Root().DisableAutoGenTag = true
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			target := filepath.Join(wd, o)
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			err = generateCodeModDocs(target)
			if err != nil {
				return err
			}
			return doc.GenMarkdownTreeCustom(cmd.Root(), o, func(_ string) string {
				return ""
			}, func(s string) string {
				return s
			})
		},
	}

	//	gendocsCmd.Flags().StringP("basepath", "b", "inventory", "Base path for the documentation (default is /inventory)")

	gendocsCmd.PreRunE = func(cmd *cobra.Command, _ []string) error {
		_ = config.BindPFlag("docs.output", cmd.Flags().Lookup("output"))
		return nil
	}
	// Define cobra flags, the default value has the lowest (least significant) precedence
	gendocsCmd.Flags().StringP("output", "o", "docs", "Output directory for the documentation (default is docs)")
	return gendocsCmd
}

func generateCodeModDocs(path string) error {
	for name, mod := range codemods.Mods {
		if mod.Description() == "" {
			continue
		}
		slog.Debug("Generating code mod docs", "name", name)
		f, err := os.Create(filepath.Join(path, name+".md"))
		if err != nil {
			return err
		}
		defer f.Close()
		f.WriteString("# " + name + " CodeMod\n\n")
		f.WriteString(mod.Description() + "\n")
		f.WriteString("\n## Usage\n\n")
		f.WriteString("```\n")
		f.WriteString(mod.Usage() + "\n")
		f.WriteString("```\n")
		slog.Debug("Generated code mod docs", "name", name)
	}
	return nil
}
