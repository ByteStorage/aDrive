package cmd

import "github.com/spf13/cobra"

func NewNameNodeServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "namenode",
		Short: "namenode operations",
	}
	cmd.AddCommand(NewNameNodeStartCmd())
	return cmd
}
