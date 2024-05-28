package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	option := TogEnvironment{
		IgnoreNewline: false,
		LogLevel:      "INFO",
		Columns:       "T,N,I,V,L,M,C,S",
		Format:        "%s    %s    %s    %s    %s    %s    %s    %s",
		Host:          "",
		ContentType:   "text",
		Locale:        "",
		TimeFormat:    "DateTime",
		Ago:           0,
	}
	envFilePath := "./tog.config"
	envFile, openErr := os.Open(envFilePath)
	if openErr != nil {
		envFilePath, err := os.UserConfigDir()
		if err == nil {
			envFilePath += "/tog.config"
			envFile, openErr = os.Open(envFilePath)

			if openErr == nil && envFile != nil {
				option = ReadEnvFile(envFile)
			} else {
				envFile, openErr = os.Create("./tog.config")
				if openErr == nil && envFile != nil {
					WriteEnvFile(envFile, option)
				} else {
					envFilePath, err = os.UserConfigDir()
					envFile, openErr = os.Create(envFilePath + "/tog.config")
					if openErr == nil && err == nil && envFile != nil {
						WriteEnvFile(envFile, option)
					}
				}
			}
		}
	} else if envFile != nil {
		option = ReadEnvFile(envFile)
	}
	envFile.Close()

	rFlag := []cli.Flag{
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
			Usage: "text print format ex) : %s	%s	%q	%s	%s	%s	%s	%s\n",
		},
		cli.StringFlag{
			Name:  "content-type, ct",
			Usage: "content type default is text. if choose json, message print json format.",
			Value: "text",
		},
		cli.StringFlag{
			Name:  "from",
			Value: time.Now().Format("2006-01-02"),
			Usage: "from date (default current day)",
		},
		cli.StringFlag{
			Name:  "to",
			Value: time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
			Usage: "to date (default tomorrow)",
		},
		cli.StringFlag{
			Name:  "log-levels, l",
			Usage: "log level list ex) INFO,ERROR",
			Value: "INFO,WARN,ERROR",
		},
		cli.StringFlag{
			Name:  "message, m",
			Usage: "message",
		},
		cli.StringFlag{
			Name:  "service-id, si",
			Usage: "service id",
		},
		cli.IntFlag{
			Name:  "tail, t",
			Usage: "tail option must be Greater than 0",
		},
		cli.BoolFlag{
			Name:  "ignore-newline, in",
			Usage: "ignore newline option (default false)",
		},
		cli.BoolFlag{
			Name:  "download, dl",
			Usage: "download option (default false)",
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
			Name:  "time-locale, tl",
			Usage: "time locale ex) Asia/Seoul | America/New_York | Europe/London ... etc. if not set, UTC",
		},
		cli.StringFlag{
			Name:  "ago",
			Usage: "time duration units for select range(current time - input_duration to current time). Valid time units are ns, us (or µs), ms, s, m, h.",
		},
	}

	rcFlag := []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Usage: "ntms-log-service host string",
		},
		cli.StringFlag{
			Name:  "from",
			Value: time.Now().Format("2006-01-02"),
			Usage: "from date (default current day)",
		},
		cli.StringFlag{
			Name:  "to",
			Value: time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
			Usage: "to date (default tomorrow)",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "rc",
			Aliases: []string{"rc"},
			Usage:   "Count(*) Gboup By Service ID In Same Service Name",
			Flags:   rcFlag,
			Action: func(c *cli.Context) error {
				if c == nil {
					return errors.New("context is nil")
				}

				hostStr := ""
				if c.IsSet("host") {
					hostStr = "http://" + c.String("host") + "/ntms-log-service/api/v1/log/amount"
				} else {
					hostStr = "http://" + option.Host + "/ntms-log-service/api/v1/log/amount"
				}
				host, err := url.Parse(hostStr)
				if err != nil {
					return err
				}

				query, err := getDefaultQueries(c, "rc", option)
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
			},
		},
		{
			Name:    "r",
			Aliases: []string{"r"},
			Usage:   "read logs ex) tog r $service-name [options...]",
			Flags:   rFlag,
			Action: func(c *cli.Context) error {
				if c == nil {
					return errors.New("context is nil")
				}

				hostStr := ""
				if c.IsSet("host") {
					hostStr = c.String("host")
				} else {
					hostStr = option.Host
				}

				if c.IsSet("tail") {
					hostStr = "ws://" + hostStr
					hostStr += "/ntms-log-service/api/v1/log/list-tail"
					err := GetLogListTail(c, hostStr, option)
					if err != nil {
						return err
					}
				} else {
					hostStr = "http://" + hostStr
					hostStr += "/ntms-log-service/api/v1/log/list"
					err := GetLogList(c, hostStr, option)
					if err != nil {
						return err
					}
				}
				return nil
			},
		},
	}

	app.Name = "tog"
	app.Usage = "TMS Log CLI"

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

type TogEnvironment struct {
	Host          string
	Ago           time.Duration
	Format        string
	Columns       string
	ContentType   string
	IgnoreNewline bool
	TimeFormat    string
	Locale        string
	LogLevel      string
}

func (t TogEnvironment) String() string {
	config := "host=" + t.Host + "\n"
	config += "ago=" + t.Ago.String() + "\n"
	config += "format=" + t.Format + "\n"
	config += "columns=" + t.Columns + "\n"
	config += "content-type=" + t.ContentType + "\n"
	config += "ignore-newline=" + strconv.FormatBool(t.IgnoreNewline) + "\n"
	config += "time-format=" + t.TimeFormat + "\n"
	config += "time-locale=" + t.Locale + "\n"
	config += "log-levels=" + t.LogLevel + "\n"
	return config
}

func ReadEnvFile(envFile *os.File) TogEnvironment {
	option := TogEnvironment{}
	scanner := bufio.NewScanner(envFile)
	for scanner.Scan() {
		temp := strings.Split(scanner.Text(), "=")
		if len(temp) == 2 && temp[1] != "" {
			opt := temp[0]
			value := temp[1]

			switch opt {
			case "host":
				option.Host = value
			case "ago":
				tempDuration, err := time.ParseDuration(value)
				if err == nil {
					option.Ago = tempDuration
				}
			case "log-levels":
				option.LogLevel = value
			case "content-type":
				option.ContentType = value
			case "ignore-newline":
				if value == "true" || value == "t" {
					option.IgnoreNewline = true
				} else {
					option.IgnoreNewline = false
				}
			case "time-format":
				option.TimeFormat = value
			case "format":
				option.Format = value
			case "columns":
				option.Columns = value
			case "time-locale":
				option.Locale = value
			}
		}
	}
	return option
}

func WriteEnvFile(envFile *os.File, option TogEnvironment) {
	config := []byte(option.String())
	envFile.Write(config)
}

func GetLogList(c *cli.Context, hostStr string, option TogEnvironment) error {
	host, err := url.Parse(hostStr)
	if err != nil {
		return err
	}

	query, err := getDefaultQueries(c, "r", option)
	if err != nil {
		return err
	}
	host.RawQuery = query.Encode()

	request := resty.New().R()
	if c.String("content-type") == "json" {
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
	if c.Bool("download") {
		file, err := os.Create(c.String("file-path") + "/" + c.String("file-name"))
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

func GetLogListTail(c *cli.Context, hostStr string, option TogEnvironment) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	host, parseErr := url.Parse(hostStr)
	if parseErr != nil {
		return parseErr
	}

	query, err := getDefaultQueries(c, "r", option)
	if err != nil {
		return err
	}
	host.RawQuery = query.Encode()

	requestHeader := make(http.Header)
	if c.String("content-type") == "json" {
		requestHeader.Add("Content-Type", "application/json")
	} else {
		requestHeader.Add("Content-Type", "text/plain")
	}
	requestHeader.Add("Content-Type", "application/json")
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
		if c.Bool("download") {
			file, err := os.Create(c.String("file-path") + "/" + c.String("file-name"))
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

func getDefaultQueries(c *cli.Context, cmd string, option TogEnvironment) (url.Values, error) {
	query := url.Values{}

	if c.IsSet("ago") {
		timeObj := time.Now()
		agoDuration, err := time.ParseDuration(c.String("ago"))
		if err != nil {
			return nil, err
		}
		query.Add("to", timeObj.Format(time.DateTime))
		query.Add("from", timeObj.Add(agoDuration*-1).Format(time.DateTime))
	} else {
		if option.Ago > 0 && !c.IsSet("from") {
			timeObj := time.Now()
			query.Add("to", timeObj.Format(time.DateTime))
			query.Add("from", timeObj.Add(option.Ago*-1).Format(time.DateTime))
		} else {
			query.Add("from", c.String("from"))
			if c.IsSet("to") {
				query.Add("to", c.String("to"))
			}
		}
	}

	query.Add("service_name", c.Args().Get(0))

	if cmd == "r" {
		if c.IsSet("log-levels") {
			levels := strings.Split(c.String("log-levels"), ",")
			for _, l := range levels {
				query.Add("log_level", l)
			}
		} else if option.LogLevel != "" {
			levels := strings.Split(option.LogLevel, ",")
			for _, l := range levels {
				query.Add("log_level", l)
			}
		}

		if c.IsSet("format") {
			query.Add("format", c.String("format"))
		} else if option.Format != "" {
			query.Add("format", option.Format)
		}
		if c.IsSet("columns") {
			query.Add("columns", c.String("columns"))
		} else if option.Columns != "" {
			query.Add("columns", option.Columns)
		}
		if c.IsSet("service-id") {
			query.Add("service_id", c.String("service-id"))
		}
		if c.IsSet("message") {
			query.Add("message", c.String("message"))
		}
		if c.IsSet("ignore-newline") {
			query.Add("ignore_newline", strconv.FormatBool(c.Bool("ignore-newline")))
		} else if option.IgnoreNewline {
			query.Add("ignore_newline", strconv.FormatBool(option.IgnoreNewline))
		}
		if c.IsSet("time-format") {
			query.Add("time_format", c.String("time-format"))
		} else if option.TimeFormat != "" {
			query.Add("time_format", option.TimeFormat)
		}
		if c.IsSet("tail") {
			query.Add("tail", strconv.Itoa(c.Int("tail")))
		}
		if c.IsSet("time-locale") {
			query.Add("time_locale", c.String("time-locale"))
		} else if option.Locale != "" {
			query.Add("time_locale", option.Locale)
		}
	}
	return query, nil
}
