package main

import (
	"os"
	"os/signal"

	"github.com/s77rt/YouTubeDownloader"
)

func main() {
	YouTubeDownloader.Init()

	server := YouTubeDownloader.NewServer()
	defer server.Close()

	ui := YouTubeDownloader.NewUI()
	defer ui.Close()

	YouTubeDownloader.BindUI(ui)

	YouTubeDownloader.LoadUI(ui, server)

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	YouTubeDownloader.Clean()
}
