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

	"github.com/bketelsen/toolbox/cobra"
	"github.com/lmittmann/tint"
	"github.com/spf13/viper"
)

var cfgFile string
var upstream string
var modsdir string
var version string
var commit string
var verbose bool

func versionString() string {
	if len(commit) > 7 {
		commit = commit[:7]
	}
	if len(commit) == 0 {
		commit = "unknown"
	}
	if len(version) == 0 {
		version = "unknown"
	}
	return fmt.Sprintf("%s (%s)", version, commit)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "surgeon",
	Version: versionString(),
	Short:   "Surgical forks of upstream repositories",

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
		slog.SetDefault(cmd.Logger)
		c, err := ReadConfig()
		if err != nil {
			cmd.Logger.Error("Reading config", tint.Err(err))
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

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .surgeon.yaml)")
	rootCmd.PersistentFlags().StringVar(&upstream, "upstream", "", "upstream repository")
	viper.BindPFlag("upstream", rootCmd.PersistentFlags().Lookup("upstream"))
	rootCmd.PersistentFlags().StringVar(&modsdir, "modsdir", "", "directory containing code modification files")
	viper.BindPFlag("modsdir", rootCmd.PersistentFlags().Lookup("modsdir"))
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "display verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	// logging

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find current directory.
		pwd, err := os.Getwd()
		cobra.CheckErr(err)

		// Search config in current directory with name ".surgeon.yaml".
		viper.AddConfigPath(pwd)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".surgeon.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		slog.Debug("config", "file", viper.ConfigFileUsed())
	}
}
