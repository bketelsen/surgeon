package main

import (
	"fmt"
	"os"

	"github.com/bketelsen/toolbox/cobra"
	mcoral "github.com/bketelsen/toolbox/mcobra"
	"github.com/muesli/roff"
	"github.com/spf13/viper"
)

// NewManCommand creates a new man command
func NewManCommand(_ *viper.Viper) *cobra.Command {
	manCmd := &cobra.Command{
		Use:                   "man",
		Short:                 "Generates surgeon's command line manpages",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		ValidArgsFunction:     cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, _ []string) error {
			manPage, err := mcoral.NewManPage(1, cmd.Root())
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
			return err
		},
	}

	return manCmd
}
