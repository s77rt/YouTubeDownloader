package YouTubeDownloader

import (
	"os"
	"path/filepath"
	E "github.com/s77rt/YouTubeDownloader/extractor"
	D "github.com/s77rt/YouTubeDownloader/downloader"
)

var (
	Version string = "dev"
)

var (
	tmpdir string
	outdir string
	extractor E.Extractor  // Currently only one is supported
	downloader *D.Downloader
	tasks Tasks
)

func Init() {
	maketmpdir()
	makeoutdir()
	extractor = E.New_YouTube_Extractor() // Currently only YouTube is supported
	downloader = D.New_Downloader()
	tasks = NewTasks()
}

func Clean() {
	os.RemoveAll(tmpdir)
}

func maketmpdir() {
	tmpdir = filepath.Join(os.TempDir(), "YouTubeDownloader")
	err := os.MkdirAll(tmpdir, 0o755)
	if err != nil {
		panic("unable to create tmpdir")
	}
}

func makeoutdir() {
	outdir = filepath.Join(".", "Downloads")
	err := os.MkdirAll(outdir, 0o755)
	if err != nil {
		panic("unable to create outdir")
	}
}
