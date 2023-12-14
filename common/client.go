package common

import (
	"github.com/xuxiaowei-com-cn/go-nexus"
	"github.com/xuxiaowei-com-cn/nexus-go/constant"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

func Client(enableLog bool, logFolder string, logName string, microseconds bool, longFile bool) (*nexus.Client, error) {
	var c = &nexus.Client{}

	flagInt := log.Ldate | log.Ltime

	if microseconds {
		flagInt = flagInt | log.Lmicroseconds
	}

	if longFile {
		flagInt = flagInt | log.Llongfile
	} else {
		flagInt = flagInt | log.Lshortfile
	}

	if logFolder == "" {
		currentUser, err := user.Current()
		if err != nil {
			return nil, err
		}
		homeDir := currentUser.HomeDir
		logFolder = filepath.Join(homeDir, constant.DefaultLogFolder)
	}

	if enableLog {
		_, err := os.Stat(logFolder)

		if os.IsNotExist(err) {
			err := os.MkdirAll(logFolder, os.ModePerm)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}

		currentTime := time.Now()

		// 格式化为日期字符串
		dateString := currentTime.Format("2006-01-02_15-04-05")

		logFile := filepath.Join(logFolder, logName+"-"+dateString+".log")

		file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}

		multi := io.MultiWriter(os.Stdout, file)

		c.Out = multi
	}

	c.Prefix = ""
	c.Flag = flagInt

	return c, nil
}
