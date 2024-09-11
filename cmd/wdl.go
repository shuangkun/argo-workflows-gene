package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

func WdlCommand() *cobra.Command {

	command := &cobra.Command{
		Use:   "wdl",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

		},
	}

	command.AddCommand(NewSubmitCommand())

	return command
}
