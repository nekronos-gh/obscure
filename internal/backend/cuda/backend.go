// Package cuda adapts the CUDA worker to the pipeline backend interface.
package cuda

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// Backend runs the CUDA worker bundled alongside the CLI.
type Backend struct{}

// Process passes the frame dimensions as arguments and streams RGB24 bytes to stdin.
func (Backend) Process(frames io.Reader, width, height int) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	worker := exec.Command(filepath.Join(filepath.Dir(executable), "worker"), strconv.Itoa(width), strconv.Itoa(height))
	worker.Stdin = frames
	worker.Stderr = os.Stderr
	if err := worker.Run(); err != nil {
		return fmt.Errorf("CUDA worker: %w", err)
	}
	return nil
}
