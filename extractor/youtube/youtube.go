package youtube

import (
	"context"
	"errors"

	"github.com/kkdai/youtube/v2"
)

type Extractor struct {
	client youtube.Client
}

func (E *Extractor) GetVideo(url string) (interface{}, error) {
	video, err := E.client.GetVideo(url)
	if err != nil {
		return nil, err
	}
	video.Formats.Sort()

	return video, nil
}

func (E *Extractor) GetVideoAndFormat(url string) (interface{}, interface{}, error) {
	video, err := E.client.GetVideo(url)
	if err != nil {
		return nil, nil, err
	}
	video.Formats.Sort()

	if len(video.Formats) == 0 {
		err = errors.New("no formats")
		return nil, nil, err
	}
	format := video.Formats[0]

	return video, format, nil
}

func (E *Extractor) GetVideos(url string) ([]interface{}, error) {
	videos := []interface{}{}

	playlist, err := E.client.GetPlaylist(url)
	if err != nil {
		video, err := E.GetVideo(url)
		if err != nil {
			return videos, err
		}
		videos = append(videos, video)
		return videos, nil
	}

	for _, v := range playlist.Videos {
		video, err := E.GetVideo(v.ID)
		if err != nil {
			return videos, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (E *Extractor) GetStreamURLContext(ctx context.Context, video interface{}, format interface{}) (string, error) {
	v, err := assertVideo(video)
	if err != nil {
		return "", err
	}
	f, err := assertFormat(format)
	if err != nil {
		return "", err
	}

	return E.client.GetStreamURLContext(ctx, v, f)
}

func (E *Extractor) GetStreamURLsContext(ctx context.Context, video interface{}, format interface{}) (string, string, error) {
	v, err := assertVideo(video)
	if err != nil {
		return "", "", err
	}
	f, err := assertFormat(format)
	if err != nil {
		return "", "", err
	}

	var url_1, url_2 string

	url_1, err = E.client.GetStreamURLContext(ctx, v, f)
	if err != nil {
		return "", "", err
	}

	if IsAudioSeparated(v, f) {
		url_2, err = E.client.GetStreamURLContext(ctx, v, &v.Formats.Type("audio")[0])
		if err != nil {
			return url_1, "", err
		}
	} else {
		url_2 = ""
	}

	return url_1, url_2, err
}

func (E *Extractor) GetFilename(video interface{}, format interface{}) (string, error) {
	v, err := assertVideo(video)
	if err != nil {
		return "", err
	}
	f, err := assertFormat(format)
	if err != nil {
		return "", err
	}

	return SanitizeFilename(v.Title) + pickIdealFileExtension(f.MimeType), nil
}

func New_Extractor() *Extractor {
	return &Extractor{youtube.Client{}}
}
