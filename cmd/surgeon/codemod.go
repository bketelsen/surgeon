/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/bketelsen/toolbox/cobra"
	"github.com/spf13/viper"
)

func NewCodemodCmd(config *viper.Viper) *cobra.Command {
	// Define our command
	codemodCmd := &cobra.Command{
		Use:   "codemod",
		Short: "Work with codemods",
		Long:  `Commands to work with codemods.`,
	}
	codemodCmd.AddCommand(NewCodemodListCmd(config))
	codemodCmd.AddCommand(NewCodemodDescribeCmd(config))
	return codemodCmd
}
