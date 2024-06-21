package togcmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
	lgrpc "tog/togcmd/grpc"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func GrpcReadLog(c *cli.Context) error {
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

	conn, err := grpc.NewClient(hostStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	client := lgrpc.NewLogClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	requestData := &lgrpc.LogDataSearchCondition{}
	if togOption.IsSet("from") {
		requestData.FromStr = *togOption.From
	} else {
		requestData.FromStr = time.Now().Format("2006-01-02")
	}

	if togOption.IsSet("to") {
		requestData.ToStr = *togOption.To
	}

	if togOption.IsSet("log-levels") {
		requestData.LogLevel = strings.Split(*togOption.LogLevels, ",")
	} else if envOption.LogLevel != "" {
		requestData.LogLevel = strings.Split(envOption.LogLevel, ",")
	}

	if togOption.IsSet("service-id") {
		requestData.ServiceId = *togOption.ServiceID
	}

	if togOption.IsSet("service-name") {
		requestData.ServiceName = *togOption.ServiceName
	}

	if togOption.IsSet("message") {
		requestData.Message = *togOption.Message
	}

	if togOption.IsSet("columns") {
		requestData.Columns = *togOption.Columns
	} else if envOption.Columns != "" {
		requestData.Columns = envOption.Columns
	}

	if togOption.IsSet("format") {
		requestData.Format = *togOption.Format
	} else if envOption.Format != "" {
		requestData.Format = envOption.Format
	}

	if togOption.IsSet("ignore-newline") {
		requestData.IgnoreNewline = *togOption.IgnoreNewline
	} else if envOption.IgnoreNewline {
		requestData.IgnoreNewline = envOption.IgnoreNewline
	}

	if togOption.IsSet("time-format") {
		requestData.TimeFormat = *togOption.TimeFormat
	} else if envOption.TimeFormat != "" {
		requestData.TimeFormat = envOption.TimeFormat
	}

	if togOption.IsSet("tail") {
		requestData.Tail = *togOption.Tail
	} else {
		requestData.Tail = "30"
	}

	if togOption.IsSet("time-zone") {
		requestData.TimeLocale = *togOption.TimeZone
	} else if envOption.TimeZone != "" {
		requestData.TimeLocale = envOption.TimeZone
	}

	r, err := client.ReadLog(ctx, requestData)
	if err != nil {
		return err
	}
	var locale *time.Location
	if requestData.TimeLocale != "" {
		locale, err = time.LoadLocation(requestData.TimeLocale)
	}

	columns := strings.Split(requestData.Columns, ",")
	tf := ParseTimeFormat(requestData)
	for {
		m, err := r.Recv()
		if err == io.EOF {
			if requestData.Tail == "" {
				r.CloseSend()
				break
			}
		} else if err != nil {
			return err
		} else {
			responseData := ""
			if togOption.IsSet("content-type") && *togOption.ContentType == "json" {
				responseData, err = JsonWithColumns(columns, locale, m)
			} else {
				if len(columns) > 0 {
					params := GetPrintArgs(columns, requestData.IgnoreNewline, tf, locale, m)
					responseData = fmt.Sprintf(
						requestData.Format,
						params...)
				} else {
					responseData = GetLogStr(requestData.IgnoreNewline, tf, locale, m)
				}
			}
			fmt.Printf(responseData + "\n")
		}
	}

	return nil
}

func JsonWithColumns(columns []string, location *time.Location, l *lgrpc.LogData) (string, error) {
	var b []byte
	var err error
	var key, value []byte

	buf := bytes.NewBuffer(b)
	buf.WriteRune('{')

	for i, c := range columns {
		switch c {
		case "T":
			if key, err = json.Marshal("reg_date"); err != nil {
				return "", err
			}
			if location == nil {
				if value, err = json.Marshal(l.RegDate); err != nil {
					return "", err
				}
			} else {
				to, err := time.Parse("2006-01-02T15:04:05.000Z", l.RegDate)
				if err != nil {
					return "", err
				} else {
					if value, err = json.Marshal(to.In(location)); err != nil {
						return "", err
					}
				}
			}
		case "N":
			if key, err = json.Marshal("service_name"); err != nil {
				return "", err
			}
			if value, err = json.Marshal(l.ServiceName); err != nil {
				return "", err
			}
		case "I":
			if key, err = json.Marshal("service_id"); err != nil {
				return "", err
			}
			if value, err = json.Marshal(l.ServiceId); err != nil {
				return "", err
			}
		case "V":
			if key, err = json.Marshal("service_version"); err != nil {
				return "", err
			}
			if value, err = json.Marshal(l.ServiceVersion); err != nil {
				return "", err
			}
		case "L":
			if key, err = json.Marshal("log_level"); err != nil {
				return "", err
			}
			if value, err = json.Marshal(l.LogLevel); err != nil {
				return "", err
			}
		case "M":
			if key, err = json.Marshal("message"); err != nil {
				return "", err
			}
			if value, err = json.Marshal(l.Message); err != nil {
				return "", err
			}
		case "C":
			if key, err = json.Marshal("caller"); err != nil {
				return "", err
			}
			if value, err = json.Marshal(l.Caller); err != nil {
				return "", err
			}
		case "S":
			if key, err = json.Marshal("stack_trace"); err != nil {
				return "", err
			}
			if value, err = json.Marshal(l.StackTrace); err != nil {
				return "", err
			}
		default:
			continue
		}
		buf.Write(key)
		buf.WriteRune(':')
		buf.Write(value)
		if i < len(columns)-1 {
			buf.WriteRune(',')
		}
	}
	buf.WriteRune('}')
	return buf.String(), nil
}

func GetPrintArgs(columns []string, isIgnoreNewline bool, timeFormat string, location *time.Location, l *lgrpc.LogData) []any {
	result := make([]any, 0)
	for _, c := range columns {
		switch c {
		case "T":
			if location == nil {
				result = append(result, l.RegDate)
			} else {
				to, err := time.Parse("2006-01-02T15:04:05.000Z", l.RegDate)
				if err != nil {
					result = append(result, l.RegDate)
				} else {
					result = append(result, to.In(location).Format(timeFormat))
				}
			}
		case "N":
			if isIgnoreNewline {
				result = append(result, strings.ReplaceAll(l.ServiceName, "\n", " "))
			} else {
				result = append(result, l.ServiceName)
			}
		case "I":
			if isIgnoreNewline {
				result = append(result, strings.ReplaceAll(l.ServiceId, "\n", " "))
			} else {
				result = append(result, l.ServiceId)
			}
		case "V":
			if isIgnoreNewline {
				result = append(result, strings.ReplaceAll(l.ServiceVersion, "\n", " "))
			} else {
				result = append(result, l.ServiceVersion)
			}
		case "L":
			if isIgnoreNewline {
				result = append(result, strings.ReplaceAll(l.LogLevel, "\n", " "))
			} else {
				result = append(result, l.LogLevel)
			}
		case "M":
			if isIgnoreNewline {
				result = append(result, strings.ReplaceAll(l.Message, "\n", " "))
			} else {
				result = append(result, l.Message)
			}
		case "C":
			if isIgnoreNewline {
				result = append(result, strings.ReplaceAll(l.Caller, "\n", " "))
			} else {
				result = append(result, l.Caller)
			}
		case "S":
			if isIgnoreNewline {
				result = append(result, strings.ReplaceAll(l.StackTrace, "\n", " "))
			} else {
				result = append(result, l.StackTrace)
			}
		}
	}
	return result
}

func GetLogStr(isIgnoreNewline bool, timeFormat string, location *time.Location, l *lgrpc.LogData) string {
	result := ""

	if location == nil {
		result += l.RegDate + "\t"
	} else {
		to, err := time.Parse("2006-01-02T15:04:05.000Z", l.RegDate)
		if err != nil {
			result = result + l.RegDate + "\t"
		} else {
			result = result + to.In(location).Format(timeFormat)
		}
	}

	if isIgnoreNewline {
		result += strings.ReplaceAll(l.ServiceName, "\n", " ") + "\t"
		result += strings.ReplaceAll(l.ServiceId, "\n", " ") + "\t"
		result += strings.ReplaceAll(l.ServiceVersion, "\n", " ") + "\t"
		result += strings.ReplaceAll(l.LogLevel, "\n", " ") + "\t"
		result += strings.ReplaceAll(l.Message, "\n", " ") + "\t"
		result += strings.ReplaceAll(l.Caller, "\n", " ") + "\t"
		result += strings.ReplaceAll(l.StackTrace, "\n", " ") + "\t"
	} else {
		result += l.ServiceName + "\t"
		result += l.ServiceId + "\t"
		result += l.ServiceVersion + "\t"
		result += l.LogLevel + "\t"
		result += l.Message + "\t"
		result += l.Caller + "\t"
		result += l.StackTrace
	}
	return result
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

func ParseTimeFormat(l *lgrpc.LogDataSearchCondition) string {
	switch l.TimeFormat {
	case "":
		return time.RFC3339
	case "RFC3339":
		return time.RFC3339
	case "RFC3339Nano":
		return time.RFC3339Nano
	case "RFC822":
		return time.RFC822
	case "RFC822Z":
		return time.RFC822Z
	case "RFC850":
		return time.RFC850
	case "RFC1123":
		return time.RFC1123
	case "RFC1123Z":
		return time.RFC1123Z
	case "ANSIC":
		return time.ANSIC
	case "DateTime":
		return time.DateTime
	case "DateOnly":
		return time.DateOnly
	}
	return l.TimeFormat
}
