package YouTubeDownloader

import (
	"os"
	"strings"
	"path/filepath"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func GetDir(path string) string {
	return filepath.Dir(path)
}

func EscapeJS(jscode string) string {
	return strings.Replace(jscode, "'", "\\'", -1)
}
