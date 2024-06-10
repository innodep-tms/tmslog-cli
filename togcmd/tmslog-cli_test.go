package togcmd

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/urfave/cli"
)

func TestInitEnvFile(t *testing.T) {
	App := cli.NewApp()

	readCom := GetReadCommand(initEnvFileCheck)

	App.Commands = []cli.Command{
		readCom,
	}
	App.Name = "tog"
	App.Usage = "TMS Log CLI"

	err := App.Run([]string{"tog", "r", "tms-dummy-service_99"})
	if err != nil {
		t.Fail()
	}
}

func initEnvFileCheck(c *cli.Context) error {
	envFile, envFilePath, envOption := InitEnvFile()
	defer envFile.Close()

	if envFile != nil {
		envFile.Close()
		file, err := os.Open(envFilePath)
		if err != nil {
			return err
		}
		fileReadData := ReadEnvFile(file)
		if fileReadData.Host != envOption.Host {
			return errors.New("host is not same")
		}
		if fileReadData.Ago != envOption.Ago {
			return errors.New("ago is not same")
		}
		if fileReadData.ContentType != envOption.ContentType {
			return errors.New("content-type is not same")
		}
		if fileReadData.Columns != envOption.Columns {
			return errors.New("download is not same")
		}
		if fileReadData.Format != envOption.Format {
			return errors.New("format is not same")
		}
		if fileReadData.IgnoreNewline != envOption.IgnoreNewline {
			return errors.New("ignore-newline is not same")
		}
		if fileReadData.LogLevel != envOption.LogLevel {
			return errors.New("log-level is not same")
		}
		if fileReadData.TimeFormat != envOption.TimeFormat {
			return errors.New("time-format is not same")
		}
		if fileReadData.TimeZone != envOption.TimeZone {
			return errors.New("time-zone is not same")
		}
	}

	return nil
}

func TestReadCommandDefaultFlagParsing(t *testing.T) {
	App := cli.NewApp()

	readCom := GetReadCommand(defaultFlagParsing)

	App.Commands = []cli.Command{
		readCom,
	}
	App.Name = "tog"
	App.Usage = "TMS Log CLI"

	err := App.Run([]string{"tog", "r", "tms-dummy-service_99", "-t", "--ct", "--from", "--to", "-l"})
	if err != nil {
		t.Fail()
	}
}

func defaultFlagParsing(c *cli.Context) error {
	togOption := ParseArgs(c)
	if !togOption.IsSet("tail") {
		return errors.New("tail flag is not set")
	} else if *togOption.Tail != "30" {
		return errors.New("tail flag is not default value 30")
	}

	if !togOption.IsSet("content-type") {
		return errors.New("content-type flag is not set")
	} else if *togOption.ContentType != "text" {
		return errors.New("content-type flag is not default value text")
	}

	if !togOption.IsSet("from") {
		return errors.New("from flag is not set")
	} else if *togOption.From != time.Now().Format("2006-01-02") {
		return errors.New("from flag is not default value " + time.Now().Format("2006-01-02"))
	}

	if !togOption.IsSet("to") {
		return errors.New("to flag is not set")
	} else if *togOption.To != time.Now().AddDate(0, 0, 1).Format("2006-01-02") {
		return errors.New("to flag is not default value " + time.Now().AddDate(0, 0, 1).Format("2006-01-02"))
	}

	if !togOption.IsSet("log-levels") {
		return errors.New("log-levels flag is not set")
	} else if *togOption.LogLevels != "INFO" {
		return errors.New("log-levels flag is not default value INFO")
	}
	return nil
}

func TestReadCommandFlagParsing(t *testing.T) {
	App := cli.NewApp()

	readCom := GetReadCommand(flagParsing)

	App.Commands = []cli.Command{
		readCom,
	}
	App.Name = "tog"
	App.Usage = "TMS Log CLI"

	err := App.Run([]string{
		"tog",
		"r",
		"tms-dummy-service_99",
		"--in=t",
		"--dl=true",
		"-t=1",
		"--host=localhost:8080",
		"--ct",
		"--from",
		"-c",
		"T,N,C",
		"--to",
		"-l=INFO,WARN,ERROR",
		"-m=hello",
		"--ago=1h"})
	if err != nil {
		t.Fail()
	}
}

func flagParsing(c *cli.Context) error {
	togOption := ParseArgs(c)
	if !togOption.IsSet("tail") {
		return errors.New("tail flag is not set")
	} else if *togOption.Tail != "1" {
		return errors.New("tail flag is not default value 30")
	}

	if !togOption.IsSet("content-type") {
		return errors.New("content-type flag is not set")
	} else if *togOption.ContentType != "text" {
		return errors.New("content-type flag is not default value text")
	}

	if !togOption.IsSet("from") {
		return errors.New("from flag is not set")
	} else if *togOption.From != time.Now().Format("2006-01-02") {
		return errors.New("from flag is not default value " + time.Now().Format("2006-01-02"))
	}

	if !togOption.IsSet("to") {
		return errors.New("to flag is not set")
	} else if *togOption.To != time.Now().AddDate(0, 0, 1).Format("2006-01-02") {
		return errors.New("to flag is not default value " + time.Now().AddDate(0, 0, 1).Format("2006-01-02"))
	}

	if !togOption.IsSet("log-levels") {
		return errors.New("log-levels flag is not set")
	} else if *togOption.LogLevels != "INFO,WARN,ERROR" {
		return errors.New("log-levels flag is not default value INFO")
	}

	if !togOption.IsSet("ignore-newline") {
		return errors.New("ignore-newline flag is not set")
	} else if *togOption.IgnoreNewline != true {
		return errors.New("ignore-newline flag is not default value true")
	}

	if !togOption.IsSet("download") {
		return errors.New("download flag is not set")
	} else if *togOption.Download != true {
		return errors.New("download flag is not default value true")
	}

	if !togOption.IsSet("host") {
		return errors.New("host flag is not set")
	} else if *togOption.Host != "localhost:8080" {
		return errors.New("host flag is not default value localhost:8080")
	}

	if !togOption.IsSet("columns") {
		return errors.New("columns flag is not set")
	} else if *togOption.Columns != "T,N,C" {
		return errors.New("columns flag is not default value T,N,C")
	}

	if !togOption.IsSet("message") {
		return errors.New("message flag is not set")
	} else if *togOption.Message != "hello" {
		return errors.New("message flag is not default value hello")
	}

	if !togOption.IsSet("ago") {
		return errors.New("ago flag is not set")
	} else if *togOption.Ago != 1*time.Hour {
		return errors.New("ago flag is not default value 1h")
	}

	if !togOption.IsSet("service-name") {
		return errors.New("service-name flag is not set")
	} else if *togOption.ServiceName != "tms-dummy-service_99" {
		return errors.New("service-name flag is not default value tms-dummy-service_99")
	}

	if togOption.IsSet("service-id") {
		return errors.New("service-id flag is wrong set")
	}

	if togOption.IsSet("time-zone") {
		return errors.New("time-zone flag is wrong set")
	}

	return nil
}

func TestReadCommandWsRequest(t *testing.T) {
	App := cli.NewApp()

	readCom := GetReadCommand(readLogTest)

	App.Commands = []cli.Command{
		readCom,
	}
	App.Name = "tog"
	App.Usage = "TMS Log CLI"

	err := App.Run([]string{
		"tog",
		"r",
		"tms-dummy-service_99",
		"--in=t",
		"--dl=true",
		"-t=1",
		"--host=localhost:8080",
		"--ct",
		"--from",
		"-c",
		"T,N,C",
		"--to",
		"-l=INFO,WARN,ERROR",
		"-m=hello",
		"--ago=1h"})
	if err != nil {
		t.Fail()
	}
}

func readLogTest(c *cli.Context) error {
	if c == nil {
		return errors.New("context is nil")
	}
	envFile, _, envOption := InitEnvFile()

	togOption := ParseArgs(c)
	if togOption.IsSet("help") {
		cli.ShowCommandHelp(c, "r")
		return nil
	}

	hostStr := ""
	if togOption.IsSet("host") {
		hostStr = *togOption.Host
		if *togOption.Host != "" && envOption.Host == "" && envFile != nil {
			WriteHostToEnvFile(envFile, *togOption.Host)
		}
	} else {
		hostStr = envOption.Host
	}
	envFile.Close()

	if togOption.IsSet("tail") {
		hostStr = "ws://" + hostStr
		hostStr += "/ntms-log-service/api/v1/log/list-tail"
		err := checkWsRequest(c, hostStr, envOption, &togOption)
		if err != nil {
			return err
		}
	} else {
		hostStr = "http://" + hostStr
		hostStr += "/ntms-log-service/api/v1/log/list"
		err := checkHttpRequest(c, hostStr, envOption, &togOption)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkWsRequest(c *cli.Context, hostStr string, envOption TogEnvironmentFile, option *TogOpt) error {
	host, parseErr := url.Parse(hostStr)
	if parseErr != nil {
		return parseErr
	}

	query, err := getReadQueries(c, envOption, option)
	if err != nil {
		return err
	}
	host.RawQuery = query.Encode()

	requestHeader := make(http.Header)
	if option.IsSet("content-type") && *option.ContentType == "json" {
		requestHeader.Add("Content-Type", "application/json")
	} else {
		requestHeader.Add("Content-Type", "text/plain")
	}

	if query.Get("to") != "" {
		return errors.New("from query is must be empty in tail")
	}

	if query.Get("tail") != "1" {
		return errors.New("tail query is must be 1 in this case")
	}

	if query.Get("service_name") != "tms-dummy-service_99" {
		return errors.New("service_name query is must be tms-dummy-service_99 in this case")
	}

	if query.Get("columns") != "T,N,C" {
		return errors.New("columns query is must be T,N,C in this case")
	}
	return nil
}

func checkHttpRequest(c *cli.Context, hostStr string, envOption TogEnvironmentFile, option *TogOpt) error {
	host, err := url.Parse(hostStr)
	if err != nil {
		return err
	}

	query, err := getReadQueries(c, envOption, option)
	if err != nil {
		return err
	}
	host.RawQuery = query.Encode()

	request := resty.New().R()
	if option.IsSet("content-type") && *option.ContentType == "json" {
		request.SetHeader("Content-Type", "application/json")
	} else {
		request.SetHeader("Content-Type", "text/plain")
	}

	fromStr := query.Get("from")
	toStr := query.Get("to")

	if toStr == "" {
		return errors.New("from query is must be not empty in this case")
	}

	if fromStr == "" {
		return errors.New("to query is must be not empty in this case")
	}

	to, err := time.Parse(time.DateTime, fromStr)
	if err != nil {
		return err
	}
	from, err := time.Parse(time.DateTime, toStr)
	if err != nil {
		return err
	}
	if from.Add(time.Hour).Format(time.DateTime) == to.Format(time.DateTime) {
		return errors.New("from and to query is must be 1 hour difference")
	}

	if query.Get("tail") != "" {
		return errors.New("tail query is must be 1 in this case")
	}

	if query.Get("service_name") != "tms-dummy-service_99" {
		return errors.New("service_name query is must be tms-dummy-service_99 in this case")
	}

	if query.Get("columns") != "T,N,C" {
		return errors.New("columns query is must be T,N,C in this case")
	}
	return nil
}

func TestReadCommandHttpRequest(t *testing.T) {
	App := cli.NewApp()

	readCom := GetReadCommand(readLogTest)

	App.Commands = []cli.Command{
		readCom,
	}
	App.Name = "tog"
	App.Usage = "TMS Log CLI"

	err := App.Run([]string{
		"tog",
		"r",
		"tms-dummy-service_99",
		"--in=t",
		"--dl=true",
		"--host=localhost:8080",
		"--ct",
		"--from",
		"-c",
		"T,N,C",
		"--to",
		"-l=INFO,WARN,ERROR",
		"-m=hello",
		"--ago=1h"})
	if err != nil {
		t.Fail()
	}
}
