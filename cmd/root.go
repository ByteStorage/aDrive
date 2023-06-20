package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aDrive",
		Short: "aDrive daemon",
	}
	cmd.AddCommand(NewNameNodeServerCmd())
	return cmd
}
