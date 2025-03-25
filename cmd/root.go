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
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/bketelsen/toolbox/cobra"
	goversion "github.com/bketelsen/toolbox/go-version"
	"github.com/bketelsen/toolbox/slug"
	"github.com/spf13/viper"
)

// var cfgFile string
var appname = "surgeon"
var cfgFile string
var (
	version   = ""
	commit    = ""
	treeState = ""
	date      = ""
	builtBy   = ""
)

var bversion = buildVersion(version, commit, date, builtBy, treeState)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "surgeon",
	Version: bversion.String(),
	InitConfig: func() *viper.Viper {
		config := viper.New()
		config.SetEnvPrefix(appname)
		config.AutomaticEnv()
		config.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", ""))
		config.SetConfigType("yaml")
		config.AddConfigPath(".") // optionally look for config in the working directory
		return config
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cfgFile != "" {
			cmd.GlobalConfig().SetConfigFile(cfgFile) // Use config file from the flag if set
		} else {
			cmd.GlobalConfig().SetConfigName(".surgeon.yaml") // name of config file

		}
		if err := cmd.GlobalConfig().ReadInConfig(); err == nil {
			slog.Info("Using config file:", slog.String("file", cmd.Config().ConfigFileUsed()))
		} else {
			slog.Error("Error reading config file", slug.Err(err))
			os.Exit(1)
		}
		// set the slog default logger to the cobra logger
		slog.SetDefault(cmd.Logger)
		// set log level based on the --verbose flag
		if cmd.GlobalConfig().GetBool("verbose") {
			cmd.SetLogLevel(slog.LevelDebug)
			cmd.Logger.Debug("Debug logging enabled")
		}
	},
	Short: "Surgical forks of upstream repositories",

	Long: `Surgeon is a tool to make surgical changes to forks of upstream repositories.

The surgeon command reads a configuration file in the current directory
named '.surgeon.yaml'.  This file contains the configuration for the
surgeon command.  The configuration file contains the upstream repository
URL, the directory containing the code modification files, and a list of
code modifications to apply to the forked repository.

The surgeon command will clone the upstream repository into a temporary directory,
then apply the code modifications to the cloned repository.  The contents of the
modified repository are copied to the current directory, overwriting any existing
files.

Important: modifications are applied in the order they are listed in the configuration,
and have a cumulative effect.  Be sure to verify your modifications before committing.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run:", cmd.GlobalConfig().AllKeys())

		c, err := ReadConfig(cmd.GlobalConfig().ConfigFileUsed())
		if err != nil {
			cmd.Logger.Error("Reading config", slug.Err(err))
			return
		}
		cmd.Logger.Debug("config", "upstream", c.Upstream, "modsdir", c.ModsDir)
		project := NewPatient(c)
		cobra.CheckErr(project.Operate())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .surgeon.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose logging")

	// logging

}

// https://www.asciiart.eu/text-to-ascii-art to make your own
// just make sure the font doesn't have backticks in the letters or
// it will break the string quoting
var asciiName = `
███████╗██╗   ██╗██████╗  ██████╗ ███████╗ ██████╗ ███╗   ██╗
██╔════╝██║   ██║██╔══██╗██╔════╝ ██╔════╝██╔═══██╗████╗  ██║
███████╗██║   ██║██████╔╝██║  ███╗█████╗  ██║   ██║██╔██╗ ██║
╚════██║██║   ██║██╔══██╗██║   ██║██╔══╝  ██║   ██║██║╚██╗██║
███████║╚██████╔╝██║  ██║╚██████╔╝███████╗╚██████╔╝██║ ╚████║
╚══════╝ ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝ ╚═════╝ ╚═╝  ╚═══╝
`

// buildVersion builds the version info for the application
func buildVersion(version, commit, date, builtBy, treeState string) goversion.Info {
	return goversion.GetVersionInfo(
		goversion.WithAppDetails(appname, "Collect and report deployment information.", "https://bketelsen.github.io/surgeon"),
		goversion.WithASCIIName(asciiName),
		func(i *goversion.Info) {
			if commit != "" {
				i.GitCommit = commit
			}
			if treeState != "" {
				i.GitTreeState = treeState
			}
			if date != "" {
				i.BuildDate = date
			}
			if version != "" {
				i.GitVersion = version
			}
			if builtBy != "" {
				i.BuiltBy = builtBy
			}

		},
	)
}
