package main

import (
	"os"
	"os/signal"
	"syscall"

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
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	select {
	case <-sigs:
	case <-ui.Done():
	}
	YouTubeDownloader.Clean()
	os.Exit(0)
}
