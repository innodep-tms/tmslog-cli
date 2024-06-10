package togcmd

import (
	"fmt"

	"github.com/urfave/cli"
)

func GetPathCommand() cli.Command {
	pathCommand := cli.Command{
		Name:    "path",
		Aliases: []string{"p"},
		Usage:   "return tog.config file path",
		Action: func(c *cli.Context) error {
			_, envFilePath, _ := InitEnvFile()
			fmt.Println(envFilePath)
			return nil
		},
	}
	return pathCommand
}
