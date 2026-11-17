// Package pipeline connects video decoding to a frame-processing backend.
package pipeline

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Backend consumes packed RGB24 frames of the given dimensions until EOF.
// Implementations must return an error if processing fails.
type Backend interface {
	Process(frames io.Reader, width, height int) error
}

// Run probes the first video stream and passes its decoded frames to backend.
func Run(path string, backend Backend) error {
	probe := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=p=0:s=x", path)
	probe.Stderr = os.Stderr
	output, err := probe.Output()
	if err != nil {
		return fmt.Errorf("inspect video: %w", err)
	}
	dimensions := strings.Split(strings.TrimSpace(string(output)), "x")
	if len(dimensions) != 2 {
		return fmt.Errorf("input has no usable video stream")
	}
	width, err := strconv.Atoi(dimensions[0])
	if err != nil || width <= 0 {
		return fmt.Errorf("invalid video width")
	}
	height, err := strconv.Atoi(dimensions[1])
	if err != nil || height <= 0 {
		return fmt.Errorf("invalid video height")
	}

	// Raw RGB24 has no frame headers. The backend uses the probed dimensions
	// to find frame boundaries, and -vsync 0 preserves the decoded frame count.
	decoder := exec.Command("ffmpeg", "-v", "error", "-i", path, "-map", "0:v:0", "-vsync", "0", "-f", "rawvideo", "-pix_fmt", "rgb24", "pipe:1")
	decoder.Stderr = os.Stderr
	frames, err := decoder.StdoutPipe()
	if err != nil {
		return fmt.Errorf("create decoder pipe: %w", err)
	}
	if err := decoder.Start(); err != nil {
		return fmt.Errorf("start decoder: %w", err)
	}
	processErr := backend.Process(frames, width, height)
	if processErr != nil {
		// The backend stopped reading, so the decoder could block writing to the pipe.
		_ = frames.Close()
		_ = decoder.Process.Kill()
	}
	decodeErr := decoder.Wait()
	if processErr != nil {
		return fmt.Errorf("process frames: %w", processErr)
	}
	if decodeErr != nil {
		return fmt.Errorf("decode video: %w", decodeErr)
	}
	return nil
}
