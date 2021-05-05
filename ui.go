package YouTubeDownloader

import (
	"fmt"
	"log"
	"net"
	"github.com/zserge/lorca"
)

func NewUI() lorca.UI {
	args := []string{}
	args = append(args, "--class=YouTube Downloader")
	ui, err := lorca.New("", "", 630, 720, args...)
	if err != nil {
		log.Fatal(err)
	}

	return ui
}

func LoadUI(ui lorca.UI, server net.Listener) {
	ui.Load(fmt.Sprintf("http://%s/frontend", server.Addr()))
}
