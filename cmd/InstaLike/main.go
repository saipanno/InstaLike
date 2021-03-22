package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mkideal/cli"
	"github.com/saipanno/go-kit/logger"

	ilike "github.com/saipanno/InstaLike/InstaLike"
	"github.com/saipanno/InstaLike/pkg/utils"
)

type args struct {
	cli.Helper

	Version bool   `cli:"version, v" usage:"verion"`
	Config  string `cli:"conf, c" usage:"specify config file" dft:"./cfg.json"`
}

func main() {

	cli.Run(new(args), func(ctx *cli.Context) (err error) {

		argv := ctx.Argv().(*args)

		if argv.Version {
			fmt.Printf("%s\n", utils.VERSION)
			return
		}

		manager := ilike.NewManager()
		if err != nil {
			return
		}

		err = manager.Start()
		if err != nil {
			return
		}

		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-exit

		logger.Info("received exit signal")
		manager.Stop()
		return
	})
}
