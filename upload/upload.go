package upload

// WallpaperUpload bing wallpaper upload interface
type WallpaperUpload interface {
	// Upload bing wallpaper Upload
	Upload(info WallpaperUploadInfo) error
}

// WallpaperUploadInfo bing wallpaper upload info
type WallpaperUploadInfo struct {
	index      uint
	market     string
	resolution string
}

// NewWallpaperUploadInfo create a new WallpaperUploadInfo
func NewWallpaperUploadInfo(index uint, market, resolution string) WallpaperUploadInfo {
	return WallpaperUploadInfo{
		index:      index,
		market:     market,
		resolution: resolution,
	}
}
