package tog

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type TogEnvironmentFile struct {
	Host          string
	Ago           time.Duration
	Format        string
	Columns       string
	ContentType   string
	IgnoreNewline bool
	TimeFormat    string
	TimeZone      string
	LogLevel      string
}

func (t TogEnvironmentFile) String() string {
	config := "ago=" + t.Ago.String() + "\n"
	config += "format=" + t.Format + "\n"
	config += "columns=" + t.Columns + "\n"
	config += "content-type=" + t.ContentType + "\n"
	config += "ignore-newline=" + strconv.FormatBool(t.IgnoreNewline) + "\n"
	config += "time-format=" + t.TimeFormat + "\n"
	config += "time-zone=" + t.TimeZone + "\n"
	config += "log-levels=" + t.LogLevel + "\n"
	return config
}

func ReadEnvFile(envFile *os.File) TogEnvironmentFile {
	option := TogEnvironmentFile{}
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
			case "time-zone":
				option.TimeZone = value
			}
		}
	}
	return option
}

func WriteEnvFile(envFile *os.File, option TogEnvironmentFile) {
	config := []byte(option.String())
	envFile.Write(config)
}

func WriteHostToEnvFile(envFile *os.File, host string) {
	config := []byte("host=" + host + "\n")
	envFile.Write(config)
}

func InitEnvFile() (*os.File, string, TogEnvironmentFile) {
	zoneName, _ := time.Now().Zone()
	envOption := TogEnvironmentFile{
		IgnoreNewline: false,
		LogLevel:      "INFO",
		Columns:       "T,N,I,V,L,M,C,S",
		Format:        "%s    %s    %s    %s    %s    %s    %s    %s",
		Host:          "",
		ContentType:   "text",
		TimeZone:      zoneName,
		TimeFormat:    "DateTime",
		Ago:           0,
	}

	var err error
	envFilePath := "./tog.config"
	envFile, openErr := os.Open(envFilePath)
	if openErr != nil {
		envFilePath, err = os.UserConfigDir()
		if err == nil {
			if runtime.GOOS == "windows" {
				envFilePath += "\\tog.config"
			} else if runtime.GOOS == "linux" {
				envFilePath += "/tog.config"
			}
			envFile, openErr = os.Open(envFilePath)
			if openErr == nil && envFile != nil {
				envOption = ReadEnvFile(envFile)
			} else {
				envFile, openErr = os.Create("./tog.config")
				if openErr == nil && envFile != nil {
					WriteEnvFile(envFile, envOption)
				} else {
					envFilePath, err = os.UserConfigDir()
					envFilePath += "/tog.config"
					envFile, openErr = os.Create(envFilePath)
					if openErr == nil && err == nil && envFile != nil {
						WriteEnvFile(envFile, envOption)
					} else {
						envFile = nil
						envFilePath = ""
					}
				}
			}
		}
	} else if envFile != nil {
		envOption = ReadEnvFile(envFile)
	}

	return envFile, envFilePath, envOption
}
