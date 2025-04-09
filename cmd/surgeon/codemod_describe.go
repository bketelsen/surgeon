/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/bketelsen/surgeon/codemods"
	"github.com/bketelsen/toolbox/cobra"
	"github.com/bketelsen/toolbox/ui"
	"github.com/spf13/viper"
)

func NewCodemodDescribeCmd(_ *viper.Viper) *cobra.Command {
	// Define our command
	describeCmd := &cobra.Command{
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

	return describeCmd
}

// describeCmd represents the describe command
