package cmd

import (
	"aDrive/api"
	"github.com/desertbit/grumble"
	"os"
)

func get(c *grumble.Context) error {
	addr := api.ApiService.NameNodeHost + ":" + api.ApiService.NameNodePort
	getResp, err := api.Get(addr, c.Args.String("file"))
	if err != nil {
		return err
	}
	err = os.WriteFile(c.Args.String("file"), getResp.Data, 0644)
	if err != nil {
		return err
	}
	return nil
}
