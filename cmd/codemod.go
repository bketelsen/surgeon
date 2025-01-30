/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// codemodCmd represents the codemod command
var codemodCmd = &cobra.Command{
	Use:   "codemod",
	Short: "Work with codemods",
	Long:  `Commands to work with codemods.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("codemod called")
	// },
}

func init() {
	rootCmd.AddCommand(codemodCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// codemodCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// codemodCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
