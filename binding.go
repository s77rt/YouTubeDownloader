package YouTubeDownloader

import (
	"github.com/zserge/lorca"
	D "github.com/s77rt/YouTubeDownloader/downloader"
)

func BindUI(ui lorca.UI) {
	ui.Bind("Ready", func() {
		ui.Eval(`setVersion('`+EscapeJS(Version)+`');`)
	})

	ui.Bind("GetVideo", func(url string) {
		video, err := extractor.GetVideo(url)
		if err != nil {
			ui.Eval(`w_GetVideo__onerror('`+EscapeJS(err.Error())+`');`)
		} else {
			ui.Eval(`add_video(`+EscapeJS(MarshalJSONS(video))+`);`)
		}
		ui.Eval(`h_unlock();`)
	})

	ui.Bind("GetVideos", func(url string) {
		videos, err := extractor.GetVideos(url)
		if err != nil {
			ui.Eval(`w_GetVideo__onerror('`+EscapeJS(err.Error())+`');`)
		} else {
			for _, video := range videos {
				ui.Eval(`add_video(`+EscapeJS(MarshalJSONS(video))+`);`)
			}
		}
		ui.Eval(`h_unlock();`)
	})

	ui.Bind("GetExistingVideo", func(uuid string, video interface{}, format interface{}) {
		filename, err := getOutputFile(extractor, downloader, video, format)
		if err != nil {
			ui.Eval(`w_GetExistingVideo__onerror('`+EscapeJS(uuid)+`', '`+EscapeJS(err.Error())+`');`)
		} else if FileExists(filename) {
			ui.Eval(`update_video_success('`+EscapeJS(uuid)+`','`+EscapeJS(filename)+`');`)
		} else {
			ui.Eval(`w_GetExistingVideo__onerror('`+EscapeJS(uuid)+`', 'not available');`)
		}
	})

	ui.Bind("DownloadVideo", func(uuid string, video interface{}, format interface{}) {
		task := tasks.New(uuid)
		filename, err := DownloadVideo(task.Context, task.Tmpdir, extractor, downloader, video, format, func(p1, p2 *D.Progress) {
			ui.Eval(`update_video_progress('`+EscapeJS(uuid)+`',`+EscapeJS(MarshalJSONS(p1))+`,`+EscapeJS(MarshalJSONS(p2))+`);`)
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
				ui.Eval(`w_DownloadVideo__onerror('`+EscapeJS(uuid)+`', '`+EscapeJS(err.Error())+`');`)
			}
		} else {
			ui.Eval(`update_video_success('`+EscapeJS(uuid)+`','`+EscapeJS(filename)+`');`)
		}
		task.Done()
	})

	ui.Bind("CancelDownload", func(uuid string) {
		err := tasks.Abort(uuid)
		if err != nil {
			ui.Eval(`w_CancelDownload__onerror('`+EscapeJS(uuid)+`', '`+EscapeJS(err.Error())+`');`)
		}
	})

	ui.Bind("PlayVideo", func(uuid string, video_filename string) {
		err := RunFile(video_filename)
		if err != nil {
			ui.Eval(`w_PlayVideo__onerror('`+EscapeJS(uuid)+`', '`+EscapeJS(err.Error())+`');`)
		}
	})

	ui.Bind("OpenFolder", func(uuid string, video_filename string) {
		err := OpenFolder(video_filename)
		if err != nil {
			ui.Eval(`w_OpenFolder__onerror('`+EscapeJS(uuid)+`', '`+EscapeJS(err.Error())+`');`)
		}
	})

	ui.Bind("DeleteVideo", func(uuid string, video_filename string) {
		err := DeleteFile(video_filename)
		if err != nil {
			ui.Eval(`w_DeleteVideo__onerror('`+EscapeJS(uuid)+`', '`+EscapeJS(err.Error())+`');`)
		}
	})
}
