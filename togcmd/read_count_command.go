package togcmd

import (
	"errors"
	"io"
	"net/url"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/urfave/cli"
)

func GetReadCountFlags() []cli.Flag {
	readCountFlag := []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Usage: "ntms-log-service host string",
		},
		cli.StringFlag{
			Name:  "from",
			Usage: "from date",
		},
		cli.StringFlag{
			Name:  "to",
			Usage: "to date",
		},
	}
	return readCountFlag
}

func GetReadCountCommand(action interface{}) cli.Command {
	readCountCommand := cli.Command{
		Name:           "rc",
		Aliases:        []string{"rc"},
		SkipArgReorder: true,
		Usage:          "Count(*) Gboup By Service ID In Same Service Name",
		Flags:          GetReadCountFlags(),
		Action:         action,
	}
	return readCountCommand
}

func ReadCountLog(c *cli.Context) error {
	if c == nil {
		return errors.New("context is nil")
	}

	envFile, _, envOption := InitEnvFile()

	hostStr := ""
	if c.IsSet("host") {
		hostStr = "http://" + c.String("host") + "/ntms-log-service/api/v1/log/amount"
		if c.String("host") != "" && envOption.Host == "" {
			WriteHostToEnvFile(envFile, c.String("host"))
		}
	} else {
		hostStr = "http://" + envOption.Host + "/ntms-log-service/api/v1/log/amount"
	}
	envFile.Close()

	host, err := url.Parse(hostStr)
	if err != nil {
		return err
	}

	query, err := getReadCountQueries(c, envOption)
	if err != nil {
		return err
	}
	host.RawQuery = query.Encode()
	resp, err := resty.New().R().SetDoNotParseResponse(true).Get(host.String())

	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return errors.New(resp.Status() + " : " + resp.String())
	}
	defer resp.RawBody().Close()

	_, err = io.Copy(os.Stdout, resp.RawBody())
	if err != nil {
		return err
	}
	return nil
}
