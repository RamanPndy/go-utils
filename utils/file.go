package goutils

import (
	"fmt"
	"io"
	"os"
)

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return err == nil, nil
}

// CopyFile copies src to dst, creating dst if it doesn't exist.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src: %w", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create dst: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	return out.Sync()
}
