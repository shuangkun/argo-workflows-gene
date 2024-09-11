package cmd

import (
	"github.com/spf13/cobra"
)

func ArgoGeneCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "argogene",
		Short: "argogene is a CLI tool for generating Argo",
		Long:  `argogene is a CLI tool for generating Argo`,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	return command
}
