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
		Usage:    "Nexus 用户名",
		Required: required,
	}
}

func PasswordFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.Password,
		Usage:    "Nexus 密码",
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
