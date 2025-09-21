package command

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ydssx/kratos-kit/cmd/cli/command/run"
	"github.com/ydssx/kratos-kit/common"
	"github.com/ydssx/kratos-kit/common/conf"
	"github.com/ydssx/kratos-kit/pkg/util"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "cli is a CLI tool for github.com/ydssx/kratos-kit",
	Long:  `cli is a CLI tool for github.com/ydssx/kratos-kit`,
}

type confKey struct{}

func WithConfig(ctx context.Context, config *conf.Bootstrap) context.Context {
	return context.WithValue(ctx, confKey{}, config)
}

func GetConfig(ctx context.Context) *conf.Bootstrap {
	return ctx.Value(confKey{}).(*conf.Bootstrap)
}

func Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		time.Sleep(time.Millisecond * 10)
	}()

	if cfgFile != "" {
		var config conf.Bootstrap
		closeConfig := conf.MustLoad(&config, cfgFile)
		defer closeConfig()

		ctx = WithConfig(ctx, &config)
	}

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "configs/config.local.yaml", "config file")

	rootCmd.AddCommand(run.RunCmd)
	rootCmd.AddCommand(NewGcsUploadCmd())
}

func NewGcsUploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gcs-upload [本地文件路径] [GCS目标路径]",
		Short: "上传文件到谷歌云存储",
		Long:  `此命令用于将本地文件上传到谷歌云存储。需要提供本地文件路径和GCS目标路径。`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			localPath := args[0]
			gcsPath := args[1]

			gcs, cleanup := common.NewGoogleCloudStorage(GetConfig(ctx))
			defer cleanup()

			file, err := os.ReadFile(localPath)
			if err != nil {
				return err
			}

			ext := filepath.Ext(localPath)
			storedName := time.Now().Format("20060102150405") + util.GenerateCode(2) + ext
			objectPath := fmt.Sprintf("%s/%s", gcsPath, storedName)

			url, err := gcs.SaveFile(ctx, "", objectPath, "", file)
			if err != nil {
				return err
			}

			cmd.Printf("文件成功上传到 %s\n", url)
			return nil
		},
	}
}
