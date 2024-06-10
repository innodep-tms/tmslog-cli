package togcmd

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
)

func GetReadFlags() []cli.Flag {
	readFlag := []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Usage: "ntms-log-service host string",
		},
		cli.StringFlag{
			Name:  "columns, c",
			Usage: "log format is T(time), N(service Name), I(service ID), V(service Version), L(log Level), M(message), C(caller), S(stack Trace). ex) T,N,I,V,L,M,C,S",
		},
		cli.StringFlag{
			Name: "format, f",
			Usage: "text print format ex) : %s	%s	%q	%s	%s	%s	%s	%s",
		},
		cli.StringFlag{
			Name:  "content-type, ct",
			Usage: "content type if choose json, message print json format.",
			Value: "text",
		},
		cli.StringFlag{
			Name:  "from",
			Value: time.Now().Format("2006-01-02"),
			Usage: "from date. this flag is ignore when ago flag is set in terminal.",
		},
		cli.StringFlag{
			Name:  "to",
			Value: time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
			Usage: "to date. this flag is ignore when ago flag is set in terminal.",
		},
		cli.StringFlag{
			Name:  "log-levels, l",
			Usage: "log level list ex) INFO,ERROR.",
			Value: "INFO",
		},
		cli.StringFlag{
			Name:  "message, m",
			Usage: "message",
		},
		cli.StringFlag{
			Name:  "service-id, si",
			Usage: "service id",
		},
		cli.StringFlag{
			Name:  "tail, t",
			Usage: "tail option.",
			Value: "30",
		},
		cli.BoolFlag{
			Name:  "ignore-newline, in",
			Usage: "ignore newline option",
		},
		cli.BoolFlag{
			Name:  "download, dl",
			Usage: "download option",
		},
		cli.StringFlag{
			Name:  "file-path, fp",
			Usage: "download path",
		},
		cli.StringFlag{
			Name:  "file-name, fn",
			Usage: "download file name",
		},
		cli.StringFlag{
			Name:  "time-format, tf",
			Usage: "golang time package's time format ex) 2006-01-02 15:04:05.000 | RFC3339 | RFC3339Nano | RFC822 | RFC822Z | RFC850 ... etc",
		},
		cli.StringFlag{
			Name:  "time-zone, tz",
			Usage: "time locale ex) Asia/Seoul | America/New_York | Europe/London ... etc. if not set, UTC",
		},
		cli.StringFlag{
			Name:  "ago",
			Usage: "time duration units for select range(current time - input_duration to current time). Valid units are ns, us, ms, s, m, h. must be greater than 0",
		},
	}
	return readFlag
}

func GetReadCommand(action interface{}) cli.Command {
	readCommand := cli.Command{
		Name:            "r",
		Aliases:         []string{"r"},
		Usage:           "read logs ex) tog r $service-name [options...]",
		SkipFlagParsing: true,
		SkipArgReorder:  true,
		Flags:           GetReadFlags(),
		Action:          action,
	}
	return readCommand
}

func ReadLog(c *cli.Context) error {
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
		err := GetLogListTail(c, hostStr, envOption, &togOption)
		if err != nil {
			return err
		}
	} else {
		hostStr = "http://" + hostStr
		hostStr += "/ntms-log-service/api/v1/log/list"
		err := GetLogList(c, hostStr, envOption, &togOption)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetLogList(c *cli.Context, hostStr string, envOption TogEnvironmentFile, option *TogOpt) error {
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

	resp, err := request.SetDoNotParseResponse(true).Get(host.String())
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return errors.New(resp.Status() + " : " + resp.String())
	}
	defer resp.RawBody().Close()

	var writer io.Writer
	if option.IsSet("download") {
		file, err := os.Create(*option.FilePath + "/" + *option.FileName)
		if err != nil {
			return err
		}
		defer file.Close()

		writer = io.MultiWriter(os.Stdout, file)
	} else {
		writer = os.Stdout
	}

	_, err = io.Copy(writer, resp.RawBody())
	if err != nil {
		return err
	}
	return nil
}

func GetLogListTail(c *cli.Context, hostStr string, envOption TogEnvironmentFile, option *TogOpt) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

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

	conn, resp, connErr := websocket.DefaultDialer.Dial(host.String(), requestHeader)
	if connErr != nil {
		return connErr
	}
	defer conn.Close()
	if resp.StatusCode != 101 {
		return errors.New(resp.Status)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		pingReceiveTimeout := time.NewTimer(5 * time.Minute)
		var writer io.Writer
		if option.IsSet("download") {
			file, err := os.Create(*option.FilePath + "/" + *option.FileName)
			if err != nil {
				return
			}
			defer file.Close()

			writer = io.MultiWriter(os.Stdout, file)
		} else {
			writer = os.Stdout
		}
		for {
			select {
			case <-pingReceiveTimeout.C:
				return
			default:
				messageType, r, err := conn.NextReader()
				if err != nil {
					return
				}
				// write message to multi writer
				if messageType == websocket.TextMessage {
					_, err := io.Copy(writer, r)
					if err != nil {
						return
					}
				}
				// close done channel when close message received
				if messageType == websocket.CloseMessage {
					return
				}
				// send pong message when ping message received
				if messageType == websocket.PingMessage {
					err := conn.WriteMessage(websocket.PongMessage, nil)
					if err != nil {
						return
					}
				}
			}
		}

	}()

	// alive check
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case t := <-ticker.C:
			// send ping message every minute
			err := conn.WriteMessage(websocket.PingMessage, []byte(t.String()))
			if err != nil {
				return err
			}
		case <-interrupt:
			// send close message when interrupt signal received
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}
