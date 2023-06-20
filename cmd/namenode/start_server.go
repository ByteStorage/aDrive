package namenode

import (
	"aDrive/daemon/namenode"
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"strconv"
)

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
