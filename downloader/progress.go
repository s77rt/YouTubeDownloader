package downloader

import (
	"time"
)

type Progress struct {
	Transferred int64
	Size int64
	Speed int64
	TimeUnit time.Duration
}
