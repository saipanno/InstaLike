package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mkideal/cli"
	"github.com/saipanno/go-kit/logger"

	ilike "github.com/saipanno/InstaLike/InstaLike"
	"github.com/saipanno/InstaLike/pkg/config"
	"github.com/saipanno/InstaLike/pkg/utils"
)

type args struct {
	cli.Helper

	Version bool   `cli:"version, v" usage:"verion"`
	Config  string `cli:"conf, c" usage:"specify config file" dft:"../../config.json"`
}

func main() {

	cli.Run(new(args), func(ctx *cli.Context) (err error) {

		argv := ctx.Argv().(*args)

		if argv.Version {
			fmt.Printf("%s\n", utils.VERSION)
			return
		}

		err = config.ParseConfig(argv.Config)
		if err != nil {
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

		err = manager.Stop()
		return
	})
}
