/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/bketelsen/surgeon/codemods"
	"github.com/bketelsen/toolbox/cobra"
	"github.com/bketelsen/toolbox/ui"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe <codemod>",
	Args:  cobra.ExactArgs(1),
	Short: "Describe a codemod",
	Long: `Describe a codemod in detail.
Show the usage and arguments for a codemod.`,
	Run: func(cmd *cobra.Command, args []string) {
		cm, ok := codemods.Mods[args[0]]
		if !ok {
			ui.Error("code mod not found", args[0])
			return
		}
		cmd.Println(cm.Usage())

	},
}

func init() {
	codemodCmd.AddCommand(describeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// describeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// describeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
