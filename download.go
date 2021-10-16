package YouTubeDownloader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	D "github.com/s77rt/YouTubeDownloader/downloader"
	E "github.com/s77rt/YouTubeDownloader/extractor"
	M "github.com/s77rt/YouTubeDownloader/merger"
)

func getOutputFile(extractor E.Extractor, downloader *D.Downloader, video interface{}, format interface{}) (string, error) {
	filename, err := extractor.GetFilename(video, format)
	if err != nil {
		return "", err
	}
	return filepath.Join(outdir, filename), nil
}

func prepareOutputFile(extractor E.Extractor, downloader *D.Downloader, video interface{}, format interface{}) (*os.File, error) {
	f, err := getOutputFile(extractor, downloader, video, format)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(f)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func DownloadVideo(ctx context.Context, tmpdir string, extractor E.Extractor, downloader *D.Downloader, video interface{}, format interface{}, progress func(*D.Progress, *D.Progress)) (string, error) {
	url_1, url_2, err := extractor.GetStreamURLsContext(ctx, video, format)
	if err != nil {
		return "", err
	}

	timeunit := time.Duration(1 * time.Second)

	var progress_1 *D.Progress
	var progress_2 *D.Progress

	if url_1 != "" && url_2 == "" {
		progress_1 = &D.Progress{TimeUnit: timeunit}
		progress_2 = progress_1

		file_1, err := downloader.DownloadFile(ctx, tmpdir, url_1, func(ctx context.Context, done chan int64, size int64, reader *D.ReaderWithCounter) {
			stop := false
			for {
				select {
				case <-ctx.Done():
					stop = true
				case <-done:
					stop = true
				default:
					progress_1.Speed = reader.Transferred - progress_1.Transferred
					progress_1.Transferred = reader.Transferred
					progress_1.Size = size
					progress(progress_1, progress_2)
				}
				if stop {
					break
				}
				time.Sleep(timeunit)
			}
		})
		if err != nil {
			return "", err
		} else {
			progress_1.Speed = 0
			progress_1.Transferred = progress_1.Size
			progress(progress_1, progress_2)
		}
		defer os.Remove(file_1.Name())

		file, err := prepareOutputFile(extractor, downloader, video, format)
		if err != nil {
			return "", err
		}
		defer file.Close()

		err = MoveFile(file_1.Name(), file.Name())
		if err != nil {
			return "", err
		}

		return file.Name(), nil

	} else if url_1 != "" && url_2 != "" {
		progress_1 = &D.Progress{TimeUnit: timeunit}
		progress_2 = &D.Progress{TimeUnit: timeunit}

		file_1, err := downloader.DownloadFile(ctx, tmpdir, url_1, func(ctx context.Context, done chan int64, size int64, reader *D.ReaderWithCounter) {
			stop := false
			for {
				select {
				case <-ctx.Done():
					stop = true
				case <-done:
					stop = true
				default:
					progress_1.Speed = reader.Transferred - progress_1.Transferred
					progress_1.Transferred = reader.Transferred
					progress_1.Size = size
					progress(progress_1, progress_2)
				}
				if stop {
					break
				}
				time.Sleep(timeunit)
			}
		})
		if err != nil {
			return "", err
		} else {
			progress_1.Speed = 0
			progress_1.Transferred = progress_1.Size
			progress(progress_1, progress_2)
		}
		defer os.Remove(file_1.Name())

		file_2, err := downloader.DownloadFile(ctx, tmpdir, url_2, func(ctx context.Context, done chan int64, size int64, reader *D.ReaderWithCounter) {
			stop := false
			for {
				select {
				case <-ctx.Done():
					stop = true
				case <-done:
					stop = true
				default:
					progress_2.Speed = reader.Transferred - progress_2.Transferred
					progress_2.Transferred = reader.Transferred
					progress_2.Size = size
					progress(progress_1, progress_2)
				}
				if stop {
					break
				}
				time.Sleep(timeunit)
			}
		})
		if err != nil {
			return "", err
		} else {
			progress_2.Speed = 0
			progress_2.Transferred = progress_2.Size
			progress(progress_1, progress_2)
		}
		defer os.Remove(file_2.Name())

		file, err := prepareOutputFile(extractor, downloader, video, format)
		if err != nil {
			return "", err
		}
		defer file.Close()

		err = M.Merge(file_1, file_2, file)
		if err != nil {
			return "", err
		}

		return file.Name(), nil
	}

	return "", fmt.Errorf("download url not found")
}
