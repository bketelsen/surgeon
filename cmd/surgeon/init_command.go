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
package main

import (
	"log/slog"
	"os"

	"github.com/bketelsen/surgeon"
	"github.com/bketelsen/toolbox/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v3"
)

func NewInitCommand(_ *viper.Viper) *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new surgical fork",
		Long: `The init command will create a new '.surgeon.yaml'
file in the current directory.	This file will contain
the configuration for the surgeon command.

Example configuration file:

	upstream: https://github.com/community-scripts/ProxmoxVE
	modsdir: mymods
`,
		Run: func(_ *cobra.Command, _ []string) {
			config := surgeon.Config{
				Upstream: "https://some.repository.com/upstream/repo",
				ModsDir:  "mymods",
				IgnoreList: []surgeon.Ignore{
					{
						Prefix: "ct",
					},
				},
				CodeMods: []surgeon.CodeMod{
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
			err = os.WriteFile(".surgeon.yaml", bb, 0o644)
			if err != nil {
				slog.Error("Writing config", "error", err)
				return
			}
		},
	}

	return initCmd
}
