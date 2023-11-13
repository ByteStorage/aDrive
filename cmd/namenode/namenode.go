package namenode

import (
	"github.com/spf13/cobra"
)

func NewNameNodeStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start namenode",
		Run:   nameNodeStart,
	}
	cmd.Flags().String("addr", "127.0.0.1:9999", "port to listen on")
	cmd.Flags().Bool("master", false, "start by boostrap")
	cmd.Flags().String("follow", "127.0.0.1:9999", "follow host")
	return cmd
}

func NewNameNodeAddServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add datanode",
		Run:   nameNodeAddServer,
	}
	cmd.Flags().String("addr", "", "namenode address")
	return cmd
}
