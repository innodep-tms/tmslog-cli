package main

import (
	"fmt"
	"os"
	"tog/togcmd"

	"github.com/urfave/cli"
)

func main() {
	App := cli.NewApp()

	pathCom := togcmd.GetPathCommand()
	readCom := togcmd.GetReadCommand(togcmd.ReadLog)
	readCountCom := togcmd.GetReadCountCommand(togcmd.ReadCountLog)

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
