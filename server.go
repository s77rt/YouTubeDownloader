package YouTubeDownloader

import (
	"embed"
	"log"
	"net"
	"net/http"
)

//go:embed frontend
var fs embed.FS

func NewServer() net.Listener {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	go http.Serve(ln, http.FileServer(http.FS(fs)))

	return ln
}
