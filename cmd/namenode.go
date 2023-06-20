package cmd

import (
	"aDrive/daemon/namenode"
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"strconv"
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

func nameNodeStart(cmd *cobra.Command, args []string) {
	addr, _ := cmd.Flags().GetString("addr")
	master, _ := cmd.Flags().GetBool("master")
	follow, _ := cmd.Flags().GetString("follow")
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		_ = fmt.Errorf("invalid address: %s", err)
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		_ = fmt.Errorf("invalid port: %s", err)
	}
	namenode.StartServer(host, master, follow, p)
}
