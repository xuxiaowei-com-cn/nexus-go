package maven

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/go-nexus"
	"github.com/xuxiaowei-com-cn/nexus-go/common"
	"github.com/xuxiaowei-com-cn/nexus-go/constant"
	"github.com/xuxiaowei-com-cn/nexus-go/flag"
	"log"
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

			flagInt, writer, err := common.Writer(enableLog, logFolder, logName, microseconds, longFile)
			if err != nil {
				return err
			}

			var c = &nexus.Client{}

			c.Flag = flagInt
			c.Out = writer

			log.SetFlags(flagInt)
			log.SetOutput(writer)

			context.App.Metadata["flag"] = flagInt
			context.App.Writer = writer

			client, err := nexus.BuildClient(c, baseUrl, username, password)
			if err != nil {
				return err
			}

			switch method {
			case "assets":
				var continuationToken = ""
				return DownloadAssets(client, repository, folder, continuationToken)
			case "browse":
				var path = ""
				return DownloadBrowse(client, repository, folder, path)
			default:

			}

			return nil
		},
	}
}

func DownloadAssets(client *nexus.Client, repository string, folder string, continuationToken string) error {

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
		err = DownloadAssets(client, repository, folder, pageAssetXO.ContinuationToken)
		if err != nil {
			return err
		}
	}

	return nil
}

func DownloadBrowse(client *nexus.Client, repository string, folder string, path string) error {

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
			err = DownloadBrowse(client, repository, folder, browse.Path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
