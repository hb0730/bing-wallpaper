package run

import (
	"errors"
	upload "github.com/TimothyYe/bing-wallpaper/upload"
	"github.com/spf13/cobra"
)

var (
	uploadType string
	resolution string
	index      uint
	market     string
)
var UploadCmd = &cobra.Command{
	Use:   "upload",
	Short: `Upload bing wallpaper to OSS`,
	Long:  `Upload bing wallpaper to OSS`,
	RunE:  runUpload,
}

func init() {
	UploadCmd.PersistentFlags().StringVar(&uploadType, "type", "s3", "The type of upload, default: `s3`")
	UploadCmd.PersistentFlags().UintVar(&index, "index", 0, "The index of wallpaper, default: `0`")
	UploadCmd.PersistentFlags().StringVar(&market, "market", "zh-CN", "The region parameter, default: `zh-CN`")
	UploadCmd.PersistentFlags().StringVar(&resolution, "resolution", "", "The resolution of wallpaper image, default: `1920`")
}

func runUpload(cmd *cobra.Command, args []string) (err error) {
	uploadInfo := upload.NewWallpaperUploadInfo(index, market, resolution)
	switch uploadType {
	case "s3":
		client, err := upload.NewS3OssClient()
		if err != nil {
			return err
		}
		return client.Upload(uploadInfo)
	default:
		err = errors.New("not support upload type")
	}
	return err
}
