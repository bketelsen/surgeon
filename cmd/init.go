/*
Copyright © 2025 Brian Ketelsen

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
package cmd

import (
	"log/slog"
	"os"

	"github.com/bketelsen/toolbox/cobra"
	yaml "gopkg.in/yaml.v3"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new surgical fork",
	Long: `The init command will create a new '.surgeon.yaml'
file in the current directory.	This file will contain
the configuration for the surgeon command.

Example configuration file:

	upstream: https://github.com/community-scripts/ProxmoxVE
	modsdir: mymods
`,
	Run: func(cmd *cobra.Command, args []string) {
		config := Config{
			Upstream: "https://some.repository.com/upstream/repo",
			ModsDir:  "mymods",
			IgnoreList: []Ignore{
				{
					Prefix: "ct",
				},
			},
			CodeMods: []CodeMod{
				{
					Description: "Modify URLS",
					Mod:         "sed",
					Match:       "cmd/*.go",
					Args:        []string{"github.com/upstream/repo", "github.com/myfork/repo"},
				},
			},
		}
		bb, err := yaml.Marshal(config)
		if err != nil {
			slog.Error("Marshal config", "error", err)
			return
		}
		err = os.WriteFile(".surgeon.yaml", bb, 0644)
		if err != nil {
			slog.Error("Writing config", "error", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
