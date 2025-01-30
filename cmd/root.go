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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var upstream string
var modsdir string
var stage bool
var commit bool
var push bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "surgeon",
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
files.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		c, err := ReadConfig()
		if err != nil {
			fmt.Println(err)
			return
		}
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
	rootCmd.PersistentFlags().BoolVar(&stage, "stage", false, "stage changes in git")
	viper.BindPFlag("stage", rootCmd.PersistentFlags().Lookup("stage"))
	rootCmd.PersistentFlags().BoolVar(&commit, "commit", false, "commit changes in git")
	viper.BindPFlag("commit", rootCmd.PersistentFlags().Lookup("commit"))
	rootCmd.PersistentFlags().BoolVar(&push, "push", false, "push changes to remote git repository")
	viper.BindPFlag("push", rootCmd.PersistentFlags().Lookup("push"))

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

		// Search config in home directory with name ".surgeon.yaml".
		viper.AddConfigPath(pwd)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".surgeon.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
