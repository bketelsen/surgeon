/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/bketelsen/toolbox/cobra"
)

// codemodCmd represents the codemod command
var codemodCmd = &cobra.Command{
	Use:   "codemod",
	Short: "Work with codemods",
	Long:  `Commands to work with codemods.`,
}

func init() {
	rootCmd.AddCommand(codemodCmd)

}
