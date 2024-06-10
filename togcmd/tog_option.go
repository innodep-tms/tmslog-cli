package togcmd

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
)

type TogOpt struct {
	setFlags      map[string]bool
	Host          *string
	Columns       *string
	Format        *string
	ContentType   *string
	From          *string
	To            *string
	LogLevels     *string
	Message       *string
	ServiceID     *string
	ServiceName   *string
	Tail          *string
	IgnoreNewline *bool
	Download      *bool
	FilePath      *string
	FileName      *string
	TimeFormat    *string
	TimeZone      *string
	Ago           *time.Duration
}

func (t TogOpt) IsSet(flagName string) bool {
	switch flagName {
	case "c":
		_, ok := t.setFlags["columns"]
		return ok
	case "f":
		_, ok := t.setFlags["format"]
		return ok
	case "ct":
		_, ok := t.setFlags["content-type"]
		return ok
	case "l":
		_, ok := t.setFlags["log-levels"]
		return ok
	case "m":
		_, ok := t.setFlags["message"]
		return ok
	case "si":
		_, ok := t.setFlags["service-id"]
		return ok
	case "t":
		_, ok := t.setFlags["tail"]
		return ok
	case "in":
		_, ok := t.setFlags["ignore-newline"]
		return ok
	case "dl":
		_, ok := t.setFlags["download"]
		return ok
	case "fp":
		_, ok := t.setFlags["file-path"]
		return ok
	case "fn":
		_, ok := t.setFlags["file-name"]
		return ok
	case "tf":
		_, ok := t.setFlags["time-format"]
		return ok
	case "tz":
		_, ok := t.setFlags["time-zone"]
		return ok
	default:
		_, ok := t.setFlags[flagName]
		return ok
	}
}

func ParseArgs(c *cli.Context) TogOpt {
	opt := TogOpt{
		setFlags: make(map[string]bool),
	}
	prevArgKind := ""
	currentArgKind := ""
	var args []string
	for _, rawArg := range c.Args() {
		if strings.Contains(rawArg, "=") {
			args = append(args, strings.Split(rawArg, "=")...)
		} else {
			args = append(args, rawArg)
			if rawArg == "--help" || rawArg == "-h" {
				opt.setFlags["help"] = true
			}
		}
	}

	for i, arg := range args {
		temp := arg
		if i == 0 {
			opt.ServiceName = &temp
			opt.setFlags["service-name"] = true
			continue
		}

		if i == 1 {
			switch arg {
			case "--host":
				prevArgKind = "host"
			case "--columns", "-c":
				prevArgKind = "columns"
			case "--format", "-f":
				prevArgKind = "format"
			case "--content-type", "--ct":
				prevArgKind = "content-type"
			case "--from":
				prevArgKind = "from"
			case "--to":
				prevArgKind = "to"
			case "--log-levels", "-l":
				prevArgKind = "log-levels"
			case "--message", "-m":
				prevArgKind = "message"
			case "--service-id", "--si":
				prevArgKind = "service-id"
			case "--tail", "-t":
				prevArgKind = "tail"
			case "--ignore-newline", "--in":
				prevArgKind = "ignore-newline"
			case "--download", "--dl":
				prevArgKind = "download"
			case "--file-path", "--fp":
				prevArgKind = "file-path"
			case "--file-name", "--fn":
				prevArgKind = "file-name"
			case "--time-format", "--tf":
				prevArgKind = "time-format"
			case "--time-zone", "--tz":
				prevArgKind = "time-zone"
			case "--ago":
				prevArgKind = "ago"
			default:
				prevArgKind = "value"
			}
			continue
		}

		switch arg {
		case "--host":
			currentArgKind = "host"
		case "--columns", "-c":
			currentArgKind = "columns"
		case "--format", "-f":
			currentArgKind = "format"
		case "--content-type", "--ct":
			currentArgKind = "content-type"
		case "--from":
			currentArgKind = "from"
		case "--to":
			currentArgKind = "to"
		case "--log-levels", "-l":
			currentArgKind = "log-levels"
		case "--message", "-m":
			currentArgKind = "message"
		case "--service-id", "--si":
			currentArgKind = "service-id"
		case "--tail", "-t":
			currentArgKind = "tail"
		case "--ignore-newline", "--in":
			currentArgKind = "ignore-newline"
		case "--download", "--dl":
			currentArgKind = "download"
		case "--file-path", "--fp":
			currentArgKind = "file-path"
		case "--file-name", "--fn":
			currentArgKind = "file-name"
		case "--time-format", "--tf":
			currentArgKind = "time-format"
		case "--time-zone", "--tz":
			currentArgKind = "time-zone"
		case "--ago":
			currentArgKind = "ago"
		default:
			currentArgKind = "value"
		}

		if prevArgKind != "value" && currentArgKind == "value" {
			switch prevArgKind {
			case "host":
				opt.setFlags["host"] = true
				opt.Host = &temp
			case "columns":
				opt.setFlags["columns"] = true
				opt.Columns = &temp
			case "format":
				opt.setFlags["format"] = true
				opt.Format = &temp
			case "content-type":
				opt.setFlags["content-type"] = true
				opt.ContentType = &temp
			case "from":
				opt.setFlags["from"] = true
				opt.From = &temp
			case "to":
				opt.setFlags["to"] = true
				opt.To = &temp
			case "log-levels":
				opt.setFlags["log-levels"] = true
				opt.LogLevels = &temp
			case "message":
				opt.setFlags["message"] = true
				opt.Message = &temp
			case "service-id":
				opt.setFlags["service-id"] = true
				opt.ServiceID = &temp
			case "tail":
				opt.setFlags["tail"] = true
				opt.Tail = &temp
			case "ignore-newline":
				if arg == "true" || arg == "t" {
					opt.setFlags["ignore-newline"] = true
					opt.IgnoreNewline = new(bool)
					*opt.IgnoreNewline = true
				} else {
					opt.setFlags["ignore-newline"] = true
					opt.IgnoreNewline = new(bool)
					*opt.IgnoreNewline = false
				}
			case "download":
				if arg == "true" || arg == "t" {
					opt.setFlags["download"] = true
					opt.Download = new(bool)
					*opt.Download = true
				} else {
					opt.setFlags["download"] = true
					opt.Download = new(bool)
					*opt.Download = false
				}
			case "file-path":
				opt.setFlags["file-path"] = true
				opt.FilePath = &temp
			case "file-name":
				opt.setFlags["file-name"] = true
				opt.FileName = &temp
			case "time-format":
				opt.setFlags["time-format"] = true
				opt.TimeFormat = &temp
			case "time-zone":
				opt.setFlags["time-zone"] = true
				opt.TimeZone = &temp
			case "ago":
				opt.Ago = new(time.Duration)
				tempDuration, err := time.ParseDuration(arg)
				if err == nil {
					opt.setFlags["ago"] = true
					*opt.Ago = tempDuration
				}
			}
		} else if prevArgKind != "value" && currentArgKind != "value" {
			switch prevArgKind {
			case "content-type":
				opt.setFlags["content-type"] = true
				defaultVal := "text"
				opt.ContentType = &defaultVal
			case "from":
				opt.setFlags["from"] = true
				defaultVal := time.Now().Format("2006-01-02")
				opt.From = &defaultVal
			case "to":
				opt.setFlags["to"] = true
				defaultVal := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
				opt.To = &defaultVal
			case "log-levels":
				opt.setFlags["log-levels"] = true
				defaultVal := "INFO"
				opt.LogLevels = &defaultVal
			case "tail":
				opt.setFlags["tail"] = true
				defaultVal := "30"
				opt.Tail = &defaultVal
			}
		}
		prevArgKind = currentArgKind
	}

	temp := args[len(args)-1]
	if prevArgKind != "value" && currentArgKind == "value" {
		switch prevArgKind {
		case "host":
			opt.setFlags["host"] = true
			opt.Host = &temp
		case "columns":
			opt.setFlags["columns"] = true
			opt.Columns = &temp
		case "format":
			opt.setFlags["format"] = true
			opt.Format = &temp
		case "content-type":
			opt.setFlags["content-type"] = true
			opt.ContentType = &temp
		case "from":
			opt.setFlags["from"] = true
			opt.From = &temp
		case "to":
			opt.setFlags["to"] = true
			opt.To = &temp
		case "log-levels":
			opt.setFlags["log-levels"] = true
			opt.LogLevels = &temp
		case "message":
			opt.setFlags["message"] = true
			opt.Message = &temp
		case "service-id":
			opt.setFlags["service-id"] = true
			opt.ServiceID = &temp
		case "tail":
			opt.setFlags["tail"] = true
			opt.Tail = &temp
		case "ignore-newline":
			if temp == "true" || temp == "t" {
				opt.setFlags["ignore-newline"] = true
				opt.IgnoreNewline = new(bool)
				*opt.IgnoreNewline = true
			} else {
				opt.setFlags["ignore-newline"] = true
				opt.IgnoreNewline = new(bool)
				*opt.IgnoreNewline = false
			}
		case "download":
			if temp == "true" || temp == "t" {
				opt.setFlags["download"] = true
				opt.Download = new(bool)
				*opt.Download = true
			} else {
				opt.setFlags["download"] = true
				opt.Download = new(bool)
				*opt.Download = false
			}
		case "file-path":
			opt.setFlags["file-path"] = true
			opt.FilePath = &temp
		case "file-name":
			opt.setFlags["file-name"] = true
			opt.FileName = &temp
		case "time-format":
			opt.setFlags["time-format"] = true
			opt.TimeFormat = &temp
		case "time-zone":
			opt.setFlags["time-zone"] = true
			opt.TimeZone = &temp
		case "ago":
			opt.Ago = new(time.Duration)
			tempDuration, err := time.ParseDuration(temp)
			if err == nil {
				opt.setFlags["ago"] = true
				*opt.Ago = tempDuration
			}
		}
	} else if prevArgKind != "value" && currentArgKind != "value" {
		switch currentArgKind {
		case "content-type":
			opt.setFlags["content-type"] = true
			defaultVal := "text"
			opt.ContentType = &defaultVal
		case "from":
			opt.setFlags["from"] = true
			defaultVal := time.Now().Format("2006-01-02")
			opt.From = &defaultVal
		case "to":
			opt.setFlags["to"] = true
			defaultVal := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
			opt.To = &defaultVal
		case "log-levels":
			opt.setFlags["log-levels"] = true
			defaultVal := "INFO"
			opt.LogLevels = &defaultVal
		case "tail":
			opt.setFlags["tail"] = true
			defaultVal := "30"
			opt.Tail = &defaultVal
		}
	}
	return opt
}

func getReadQueries(c *cli.Context, envOption TogEnvironmentFile, option *TogOpt) (url.Values, error) {
	query := url.Values{}

	if option != nil {
		if option.IsSet("ago") {
			timeObj := time.Now()
			agoDuration, err := time.ParseDuration(option.Ago.String())
			if err != nil {
				return nil, err
			}
			query.Add("to", timeObj.Format(time.DateTime))
			query.Add("from", timeObj.Add(agoDuration*-1).Format(time.DateTime))
		} else {
			if envOption.Ago > 0 && !c.IsSet("from") {
				timeObj := time.Now()
				query.Add("to", timeObj.Format(time.DateTime))
				query.Add("from", timeObj.Add(envOption.Ago*-1).Format(time.DateTime))
			} else {
				query.Add("from", c.String("from"))
				if c.IsSet("to") {
					query.Add("to", c.String("to"))
				}
			}
		}
		query.Add("service_name", *option.ServiceName)
	} else {
		if c.IsSet("ago") {
			timeObj := time.Now()
			agoDuration, err := time.ParseDuration(c.String("ago"))
			if err != nil {
				return nil, err
			}
			query.Add("to", timeObj.Format(time.DateTime))
			query.Add("from", timeObj.Add(agoDuration*-1).Format(time.DateTime))
			query.Add("service_name", *option.Host)
		} else {
			if envOption.Ago > 0 && !c.IsSet("from") {
				timeObj := time.Now()
				query.Add("to", timeObj.Format(time.DateTime))
				query.Add("from", timeObj.Add(envOption.Ago*-1).Format(time.DateTime))
			} else {
				query.Add("from", c.String("from"))
				if c.IsSet("to") {
					query.Add("to", c.String("to"))
				}
			}
			query.Add("service_name", c.Args().Get(0))
		}
	}

	if option.IsSet("log-levels") {
		levels := GetLogLevelList(*option.LogLevels)
		for _, l := range levels {
			query.Add("log_level", l)
		}
	} else if envOption.LogLevel != "" {
		levels := GetLogLevelList(envOption.LogLevel)
		for _, l := range levels {

			query.Add("log_level", l)
		}
	} else {
		query.Add("log_level", "INFO")
	}

	if option.IsSet("format") {
		query.Add("format", *option.Format)
	} else if envOption.Format != "" {
		query.Add("format", envOption.Format)
	}
	if option.IsSet("columns") {
		query.Add("columns", *option.Columns)
	} else if envOption.Columns != "" {
		query.Add("columns", envOption.Columns)
	}
	if option.IsSet("service-id") {
		query.Add("service_id", *option.ServiceID)
	}
	if option.IsSet("message") {
		query.Add("message", *option.Message)
	}
	if option.IsSet("ignore-newline") {
		query.Add("ignore_newline", strconv.FormatBool(*option.IgnoreNewline))
	} else if envOption.IgnoreNewline {
		query.Add("ignore_newline", strconv.FormatBool(envOption.IgnoreNewline))
	}
	if option.IsSet("time-format") {
		query.Add("time_format", *option.TimeFormat)
	} else if envOption.TimeFormat != "" {
		query.Add("time_format", envOption.TimeFormat)
	}
	if option.IsSet("tail") {
		if _, err := strconv.Atoi(*option.Tail); err == nil {
			query.Add("tail", *option.Tail)
		} else {
			query.Add("tail", "30")
		}
	}
	if option.IsSet("time-zone") {
		query.Add("time_locale", *option.TimeZone)
	} else if envOption.TimeZone != "" {
		query.Add("time_locale", envOption.TimeZone)
	}
	return query, nil
}

func getReadCountQueries(c *cli.Context, envOption TogEnvironmentFile) (url.Values, error) {
	query := url.Values{}

	if c.IsSet("from") {
		query.Add("from", c.String("from"))
	}
	if c.IsSet("to") {
		query.Add("to", c.String("to"))
	}
	query.Add("service_name", c.Args().Get(0))
	return query, nil
}

func GetLogLevelList(rawString string) []string {
	result := []string{}
	for _, rawLevel := range strings.Split(rawString, ",") {
		rawLevel = strings.ToUpper(rawLevel)
		switch rawLevel {
		case "D":
			result = append(result, "DEBUG")
		case "I":
			result = append(result, "INFO")
		case "W":
			result = append(result, "WARN")
		case "E":
			result = append(result, "ERROR")
		case "F":
			result = append(result, "FATAL")
		case "P":
			result = append(result, "PANIC")
		default:
			result = append(result, rawLevel)
		}
	}
	return result
}
