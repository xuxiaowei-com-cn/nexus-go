package maven

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/go-nexus"
	"github.com/xuxiaowei-com-cn/nexus-go/common"
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
			flag.MethodFlag(true), flag.FolderFlag(true),
			flag.EnableLogFlag(), flag.LogNameFlag(), flag.LogFolderFlag(), flag.MicrosecondsFlag(), flag.LongFileFlag()),
		Action: func(context *cli.Context) error {
			var baseUrl = context.String(constant.BaseUrl)
			var username = context.String(constant.Username)
			var password = context.String(constant.Password)
			var method = context.String(constant.Method)
			var repository = context.String(constant.Repository)
			var folder = context.String(constant.Folder)
			var enableLog = context.Bool(constant.EnableLog)
			var logName = context.String(constant.LogName)
			var logFolder = context.String(constant.LogFolder)
			var microseconds = context.Bool(constant.Microseconds)
			var longFile = context.Bool(constant.LongFile)

			c, err := common.Client(enableLog, logFolder, logName, microseconds, longFile)
			if err != nil {
				return err
			}

			switch method {
			case "assets":
				var continuationToken = ""
				return DownloadAssets(baseUrl, username, password, repository, folder, continuationToken, c)
			case "browse":
				var path = ""
				return DownloadBrowse(baseUrl, username, password, repository, folder, path, c)
			default:

			}

			return nil
		},
	}
}

func DownloadAssets(baseUrl string, username string, password string, repository string, folder string, continuationToken string,
	c *nexus.Client) error {

	client, err := nexus.BuildClient(c, baseUrl, username, password)
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
		err = DownloadAssets(baseUrl, username, password, repository, folder, pageAssetXO.ContinuationToken, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func DownloadBrowse(baseUrl string, username string, password string, repository string, folder string, path string,
	c *nexus.Client) error {

	client, err := nexus.BuildClient(c, baseUrl, username, password)
	if err != nil {
		return err
	}

	browses, response, err := client.ExtDirect.GetBrowseRepository(repository, path)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("GetBrowseRepository status %s", response.Status)
	}

	for _, browse := range browses {
		if browse.Type == "file" {
			var filePath = filepath.Join(folder, browse.Path)
			var fileFolder = filepath.Dir(filePath)

			err = os.MkdirAll(fileFolder, os.ModePerm)
			if err != nil {
				return err
			}

			response, err = client.File.Download(http.MethodGet, browse.Url, filePath, nil, nil)
			if err != nil {
				return err
			}
			if response.StatusCode != http.StatusOK {
				return fmt.Errorf("Download Browse status %s ", response.Status)
			}
		} else {
			err = DownloadBrowse(baseUrl, username, password, repository, folder, browse.Path, c)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
