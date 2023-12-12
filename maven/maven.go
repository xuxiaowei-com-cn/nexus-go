package maven

import (
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/nexus-go/flag"
)

func MavenCommand() *cli.Command {
	return &cli.Command{
		Name:  "maven",
		Usage: "Maven 仓库",
		Flags: append(flag.Common(false), flag.RepositoryFlag(false), flag.MethodFlag(false),
			flag.FolderFlag(false)),
		Subcommands: []*cli.Command{
			DownloadCommand(),
		},
	}
}
