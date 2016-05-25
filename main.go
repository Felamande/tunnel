package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/qiniu/log"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetPrefix("[tunnel]")
	log.SetFlags(log.LdefaultShort)
}

func main() {
	app := cli.NewApp()
	app.Commands = append(app.Commands, serverCmd(), clientCmd())
	log.Error(app.Run(os.Args))
}
