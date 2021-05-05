package downloader

import (
	"io"
	"os"
	"fmt"
	"net"
	"net/http"
	"io/ioutil"
	"strconv"
	"time"
	"context"
)

type Downloader struct {
	client *http.Client
}

func (D *Downloader) DownloadFile(ctx context.Context, tmpdir string, url string, callback func(context.Context, chan int64, int64, *ReaderWithCounter)) (*os.File, error) {
	// Create a temporarily file
	file, err := ioutil.TempFile(tmpdir, "file.*.tmp")
	if err != nil {
		return nil, err
	}

	// Get header info
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Range", "bytes=0-")
	resp, err := D.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Get file size
	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return nil, err
	}

	// Callback
	done := make(chan int64)
	reader := &ReaderWithCounter{}
	go callback(ctx, done, int64(size), reader)

	// Get the data
	req, err = http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		done <- 0
		return nil, err
	}
	req.Header.Set("Range", "bytes=0-")
	resp, err = D.client.Do(req)
	if err != nil {
		done <- 0
		return nil, err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		done <- 0
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	reader.Reader = resp.Body
	written, err := io.Copy(file, reader)

	if err != nil {
		done <- 0
		return nil, err
	}

	done <- written
	return file, nil
}

func New_Downloader() *Downloader {
	httpTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	return &Downloader{
		&http.Client{Transport: httpTransport},
	}
}
