package YouTubeDownloader

import (
	"fmt"
	"os"
	"github.com/skratchdot/open-golang/open"
)

func RunFile(file string) error {
	if !FileExists(file) {
		return fmt.Errorf("file deos not exist")
	}

	return open.Run(file)
}

func OpenFolder(file string) error {
	if !FileExists(file) {
		return fmt.Errorf("file deos not exist")
	}

	dir := GetDir(file)

	return open.Run(dir)
}

func DeleteFile(file string) error {
	if !FileExists(file) {
		return fmt.Errorf("file deos not exist")
	}

	return os.Remove(file)
}
