package main

import (
	"fmt"
	"os"
	"tog/tog"

	"github.com/urfave/cli"
)

func main() {
	App := cli.NewApp()

	pathCom := tog.GetPathCommand()
	readCom := tog.GetReadCommand(tog.ReadLog)
	readCountCom := tog.GetReadCountCommand(tog.ReadCountLog)

	App.Commands = []cli.Command{
		readCountCom,
		readCom,
		pathCom,
	}
	App.Name = "tog"
	App.Usage = "TMS Log CLI"

	err := App.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
