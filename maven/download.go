package maven

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/go-nexus"
	"github.com/xuxiaowei-com-cn/nexus-go/constant"
	"github.com/xuxiaowei-com-cn/nexus-go/flag"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadCommand() *cli.Command {
	return &cli.Command{
		Name:    "download",
		Aliases: []string{"dl"},
		Usage:   "下载",
		Flags: append(flag.Common(true), flag.RepositoryFlag(true),
			flag.MethodFlag(true), flag.FolderFlag(true)),
		Action: func(context *cli.Context) error {
			var baseUrl = context.String(constant.BaseUrl)
			var username = context.String(constant.Username)
			var password = context.String(constant.Password)
			var method = context.String(constant.Method)
			var repository = context.String(constant.Repository)
			var folder = context.String(constant.Folder)

			switch method {
			case "assets":
				var continuationToken = ""
				return DownloadAssets(baseUrl, username, password, repository, folder, continuationToken)
			default:

			}

			return nil
		},
	}
}

func DownloadAssets(baseUrl string, username string, password string, repository string, folder, continuationToken string) error {

	client, err := nexus.NewClient(baseUrl, username, password)
	if err != nil {
		return err
	}

	requestQuery := &nexus.ListAssetsQuery{
		Repository: repository,
	}

	if continuationToken != "" {
		requestQuery.ContinuationToken = continuationToken
	}

	pageAssetXO, response, err := client.Assets.ListAssets(requestQuery)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("ListAssets status %s (%s)", response.Status, continuationToken)
	}

	for _, assetXO := range pageAssetXO.Items {

		var filePath = filepath.Join(folder, assetXO.Path)
		var fileFolder = filepath.Dir(filePath)

		err = os.MkdirAll(fileFolder, os.ModePerm)
		if err != nil {
			return err
		}

		response, err = client.File.Download(http.MethodGet, assetXO.DownloadUrl, filePath, nil, nil)
		if err != nil {
			return err
		}
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("Download assets status %s (%s) ", response.Status, continuationToken)
		}
	}

	if pageAssetXO.ContinuationToken != "" {
		err = DownloadAssets(baseUrl, username, password, repository, folder, pageAssetXO.ContinuationToken)
		if err != nil {
			return err
		}
	}

	return nil
}
