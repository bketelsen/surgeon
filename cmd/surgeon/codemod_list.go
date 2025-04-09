/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/bketelsen/surgeon/codemods"
	"github.com/spf13/viper"

	"github.com/bketelsen/toolbox/cobra"
	"github.com/bketelsen/toolbox/ui"
)

func NewCodemodListCmd(_ *viper.Viper) *cobra.Command {
	// Define our command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available codemods",
		Long:  `List all available codemods.`,
		Run: func(cmd *cobra.Command, _ []string) {
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

	return listCmd
}

type List struct {
	Name        string `table:"name,default_sort"`
	Description string `table:"description"`
}
