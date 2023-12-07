package cmd

import (
	"aDrive/api"
	"fmt"
	"github.com/desertbit/grumble"
	"os"
)

func put(c *grumble.Context) error {
	addr := api.ApiService.NameNodeHost + ":" + api.ApiService.NameNodePort
	bytes, err := os.ReadFile(c.Args.String("file"))
	if err != nil {
		return err
	}
	putResp, err := api.Put(addr, c.Args.String("directory"), bytes)
	if err != nil {
		return err
	}
	fmt.Println(putResp.DataMessage)
	return nil
}
