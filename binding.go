package YouTubeDownloader

import (
	"github.com/zserge/lorca"
	D "github.com/s77rt/YouTubeDownloader/downloader"
	M "github.com/s77rt/YouTubeDownloader/merger"
)

func BindUI(ui lorca.UI) {
	ui.Bind("Ready", func() {
		ui.Eval(`setVersion(`+MarshalJSONS(Version)+`);`)
		ui.Eval(`w_CheckFFmpeg();`)
	})

	ui.Bind("CheckFFmpeg", func() {
		if exists := M.FFmpeg_exists(); !exists {
			ui.Eval(`w_CheckFFmpeg__onerror("FFmpeg is not installed");`)
		}
	})

	ui.Bind("GetVideo", func(url string) {
		video, err := extractor.GetVideo(url)
		if err != nil {
			ui.Eval(`w_GetVideo__onerror(`+MarshalJSONS(err.Error())+`);`)
		} else {
			ui.Eval(`add_video(`+MarshalJSONS(video)+`);`)
		}
		ui.Eval(`h_unlock();`)
	})

	ui.Bind("GetVideos", func(url string) {
		videos, err := extractor.GetVideos(url)
		if err != nil {
			ui.Eval(`w_GetVideo__onerror(`+MarshalJSONS(err.Error())+`);`)
		} else {
			for _, video := range videos {
				ui.Eval(`add_video(`+MarshalJSONS(video)+`);`)
			}
		}
		ui.Eval(`h_unlock();`)
	})

	ui.Bind("GetExistingVideo", func(uuid string, video interface{}, formats []interface{}) {
		found := false
		for _, format := range(formats) {
			filename, err := getOutputFile(extractor, downloader, video, format)
			if err != nil {
				ui.Eval(`w_GetExistingVideo__onerror(`+MarshalJSONS(uuid)+`, `+MarshalJSONS(err.Error())+`);`)
				break
			}
			if FileExists(filename) {
				found = true
				ui.Eval(`update_video_success(`+MarshalJSONS(uuid)+`,`+MarshalJSONS(filename)+`);`)
				break
			}
		}
		if !found {
			ui.Eval(`w_GetExistingVideo__onerror(`+MarshalJSONS(uuid)+`, "not available");`)
		}
	})

	ui.Bind("DownloadVideo", func(uuid string, video interface{}, format interface{}) {
		task := tasks.New(uuid)
		filename, err := DownloadVideo(task.Context, task.Tmpdir, extractor, downloader, video, format, func(p1, p2 *D.Progress) {
			ui.Eval(`update_video_progress(`+MarshalJSONS(uuid)+`,`+MarshalJSONS(p1)+`,`+MarshalJSONS(p2)+`);`)
		})
		if err != nil {
			if task.Context.Err() != nil {
				/*
				Errors that are catched in this state where the context is closed
				are probably fine to ignore, because as the context get closed, by default a "post" function
				get executed and will delete the tmp files for that task, which may result in "no such file or directory" error
				for a function that may still trying to use such temp files, thus the error.
				the error may also be context.Canceled which is fine.
				this should be further investigated...
				*/
			} else {
				ui.Eval(`w_DownloadVideo__onerror(`+MarshalJSONS(uuid)+`, `+MarshalJSONS(err.Error())+`);`)
			}
		} else {
			ui.Eval(`update_video_success(`+MarshalJSONS(uuid)+`,`+MarshalJSONS(filename)+`);`)
		}
		task.Done()
	})

	ui.Bind("CancelDownload", func(uuid string) {
		err := tasks.Abort(uuid)
		if err != nil {
			ui.Eval(`w_CancelDownload__onerror(`+MarshalJSONS(uuid)+`, `+MarshalJSONS(err.Error())+`);`)
		}
	})

	ui.Bind("PlayVideo", func(uuid string, video_filename string) {
		err := RunFile(video_filename)
		if err != nil {
			ui.Eval(`w_PlayVideo__onerror(`+MarshalJSONS(uuid)+`, `+MarshalJSONS(err.Error())+`);`)
		}
	})

	ui.Bind("OpenFolder", func(uuid string, video_filename string) {
		err := OpenFolder(video_filename)
		if err != nil {
			ui.Eval(`w_OpenFolder__onerror(`+MarshalJSONS(uuid)+`, `+MarshalJSONS(err.Error())+`);`)
		}
	})

	ui.Bind("DeleteVideo", func(uuid string, video_filename string) {
		err := DeleteFile(video_filename)
		if err != nil {
			ui.Eval(`w_DeleteVideo__onerror(`+MarshalJSONS(uuid)+`, `+MarshalJSONS(err.Error())+`);`)
		}
	})
}
