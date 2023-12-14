package common

import (
	"github.com/xuxiaowei-com-cn/nexus-go/constant"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

func Writer(enableLog bool, logFolder string, logName string, microseconds bool, longFile bool) (int, io.Writer, error) {

	flag := log.Ldate | log.Ltime

	if microseconds {
		flag = flag | log.Lmicroseconds
	}

	if longFile {
		flag = flag | log.Llongfile
	} else {
		flag = flag | log.Lshortfile
	}

	if logFolder == "" {
		currentUser, err := user.Current()
		if err != nil {
			return flag, nil, err
		}
		homeDir := currentUser.HomeDir
		logFolder = filepath.Join(homeDir, constant.DefaultLogFolder)
	}

	if enableLog {
		_, err := os.Stat(logFolder)

		if os.IsNotExist(err) {
			err := os.MkdirAll(logFolder, os.ModePerm)
			if err != nil {
				return flag, nil, err
			}
		} else if err != nil {
			return flag, nil, err
		}

		currentTime := time.Now()

		// 格式化为日期字符串
		dateString := currentTime.Format("2006-01-02_15-04-05")

		logFile := filepath.Join(logFolder, logName+"-"+dateString+".log")

		file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return flag, nil, err
		}

		multi := io.MultiWriter(os.Stdout, file)

		return flag, multi, nil
	}

	return flag, nil, nil
}
