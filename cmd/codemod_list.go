/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/bketelsen/surgeon/codemods"

	"github.com/bketelsen/toolbox/ui"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available codemods",
	Long:  `List all available codemods.`,
	Run: func(cmd *cobra.Command, args []string) {
		var list []List
		// Iterate over the mods and add them to the list
		for name, mod := range codemods.Mods {
			list = append(list, List{
				Name:        name,
				Description: mod.Description(),
			})
		}
		// Print the list
		out, err := ui.DisplayTable(list, "", nil)
		cobra.CheckErr(err)
		cmd.Println(out)

	},
}

func init() {
	codemodCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type List struct {
	Name        string `table:"name,default_sort"`
	Description string `table:"description"`
}
