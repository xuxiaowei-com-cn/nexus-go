package maven

import (
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/go-nexus"
	"github.com/xuxiaowei-com-cn/nexus-go/constant"
	"github.com/xuxiaowei-com-cn/nexus-go/flag"
)

func UploadCommand() *cli.Command {
	return &cli.Command{
		Name:    "upload",
		Aliases: []string{"up"},
		Usage:   "上传",
		Flags:   append(flag.Common(true), flag.RepositoryFlag(true), flag.FolderFlag(true)),
		Action: func(context *cli.Context) error {
			var baseUrl = context.String(constant.BaseUrl)
			var username = context.String(constant.Username)
			var password = context.String(constant.Password)
			var repository = context.String(constant.Repository)
			var folder = context.String(constant.Folder)

			return Upload(baseUrl, username, password, repository, folder)
		},
	}
}

func Upload(baseUrl string, username string, password string, repositoryName string, folder string) error {

	client, err := nexus.NewClient(baseUrl, username, password)
	if err != nil {
		return err
	}

	return client.Repository.UploadFolder(folder, repositoryName)
}
