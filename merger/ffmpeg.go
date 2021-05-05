package merger

import (
	"os"
	"os/exec"
	"fmt"
)

func ffmpeg_merge(videoFile, audioFile, destFile *os.File) error {
	ffmpegVersionCmd := exec.Command("ffmpeg", "-y",
		"-i", videoFile.Name(),
		"-i", audioFile.Name(),
		"-c", "copy", // copy without re-encoding
		"-shortest", // Finish encoding when the shortest input stream ends
		destFile.Name(),
		"-loglevel", "warning",
	)
	ffmpegVersionCmd.Stderr = os.Stderr
	ffmpegVersionCmd.Stdout = os.Stdout

	err := ffmpegVersionCmd.Run()
	if err != nil {
		os.Remove(destFile.Name())
		return fmt.Errorf("FFmpeg: %s", err.Error())
	}
	return nil
}

func ffmpeg_merge_fallback(videoFile, audioFile, destFile *os.File) error {
	ffmpegVersionCmd := exec.Command("ffmpeg", "-y",
		"-i", videoFile.Name(),
		"-i", audioFile.Name(),
		"-c:v", "copy", // copy without re-encoding
		"-c:a", "libvorbis", // libvorbis re-encoding
		"-shortest", // Finish encoding when the shortest input stream ends
		destFile.Name(),
		"-loglevel", "warning",
	)
	ffmpegVersionCmd.Stderr = os.Stderr
	ffmpegVersionCmd.Stdout = os.Stdout

	err := ffmpegVersionCmd.Run()
	if err != nil {
		os.Remove(destFile.Name())
		return fmt.Errorf("FFmpeg: %s", err.Error())
	}
	return nil
}

func Merge(videoFile, audioFile, destFile *os.File) error {
	err := ffmpeg_merge(videoFile, audioFile, destFile)
	if err != nil {
		err = ffmpeg_merge_fallback(videoFile, audioFile, destFile)
	}
	return err
}
