package cmd

import (
	"aDrive/cmd/namenode"
	"github.com/spf13/cobra"
)

func NewNameNodeServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "namenode",
		Short: "namenode operations",
	}
	cmd.AddCommand(namenode.NewNameNodeStartCmd())
	return cmd
}
