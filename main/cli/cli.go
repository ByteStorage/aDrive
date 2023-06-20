package cli

import (
	"aDrive/daemon/datanode"
	"aDrive/daemon/namenode"
	"flag"
	"log"
	"os"
)

func StartServer() {
	dataNodeCommand := flag.NewFlagSet("datanode", flag.ExitOnError)
	nameNodeCommand := flag.NewFlagSet("namenode", flag.ExitOnError)
	clientCommand := flag.NewFlagSet("client", flag.ExitOnError)

	nameNodeAddr := dataNodeCommand.String("namenode", "116.62.156.91:9999", "NameNode communication port")
	dataport := dataNodeCommand.Int("port", 7000, "")
	dataLocation := dataNodeCommand.String("path", "data/", "")
	datanodeHost := dataNodeCommand.String("host", "", "")

	master := nameNodeCommand.Bool("master", false, "start by boostrap")
	follow := nameNodeCommand.String("follow", "", "")
	port := nameNodeCommand.Int("port", 9999, "")
	host := nameNodeCommand.String("host", "", "")

	if len(os.Args) < 2 {
		log.Println("sub-command is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "datanode":
		_ = dataNodeCommand.Parse(os.Args[2:])
		datanode.StartServer(*datanodeHost, *nameNodeAddr, *dataport, *dataLocation)

	case "namenode":
		_ = nameNodeCommand.Parse(os.Args[2:])
		namenode.StartServer(*host, *master, *follow, *port)

	case "client":
		_ = clientCommand.Parse(os.Args[2:])

	}
}
