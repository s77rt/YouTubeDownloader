package extractor

import (
	"context"

	"github.com/s77rt/YouTubeDownloader/extractor/youtube"
)

type Extractor interface {
	GetVideo(url string) (interface{}, error)
	GetVideoAndFormat(url string) (interface{}, interface{}, error)
	GetVideos(url string) ([]interface{}, error)
	GetStreamURLContext(ctx context.Context, video interface{}, format interface{}) (string, error)
	GetStreamURLsContext(ctx context.Context, video interface{}, format interface{}) (string, string, error)
	GetFilename(video interface{}, format interface{}) (string, error)
}

func New_YouTube_Extractor() Extractor {
	return youtube.New_Extractor()
}
