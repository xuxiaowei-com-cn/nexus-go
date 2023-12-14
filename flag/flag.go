package flag

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/nexus-go/constant"
)

func BaseUrlFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.BaseUrl,
		Usage:    "Nexus URL",
		Required: required,
	}
}

func UsernameFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.Username,
		Usage:    "Nexus 用户名，匿名访问请填写空",
		Required: required,
	}
}

func PasswordFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.Password,
		Usage:    "Nexus 密码，匿名访问请填写空",
		Required: required,
	}
}

func Common(required bool) []cli.Flag {
	return []cli.Flag{
		BaseUrlFlag(required),
		UsernameFlag(required),
		PasswordFlag(required),
	}
}

func MethodFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.Method,
		Usage:    "下载方法，如：assets、browse",
		Value:    "assets",
		Required: required,
		Action: func(context *cli.Context, s string) error {

			var method = context.String(constant.Method)

			switch method {
			case "assets":
				break
			case "browse":
				break
			default:
				return fmt.Errorf("method 错误，请输入：assets、browse")
			}

			return nil
		},
	}
}

func RepositoryFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.Repository,
		Usage:    "仓库名称",
		Required: required,
	}
}

func FolderFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.Folder,
		Usage:    "文件夹",
		Required: required,
	}
}

func EnableLogFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:  constant.EnableLog,
		Usage: "开启日志",
		Value: false,
	}
}

func LogNameFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  constant.LogName,
		Usage: "日志名称-前缀",
		Value: constant.DefaultLogName,
	}
}

func LogFolderFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  constant.LogFolder,
		Usage: fmt.Sprintf("日志文件夹，默认是当前用户主目录下的 %s 文件夹", constant.DefaultLogFolder),
	}
}

func MicrosecondsFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:  constant.Microseconds,
		Usage: "日志打印时间精确到微秒",
		Value: false,
	}
}

func LongFileFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:  constant.LongFile,
		Usage: "日志打印使用长包名",
		Value: false,
	}
}
