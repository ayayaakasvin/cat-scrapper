// Package provides ready functions which will simplify action of fetching and saving images of cat via libs
// github.com/ayayaakasvin/cat-photo-fetch
package fsengine

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ayayaakasvin/cat-scrapper/internal/domain"
)

type FSEngine struct {
	savePath string
}

func NewFSE(savePath string) (domain.ImageFileSystem, error) {
	if savePath == "" {
		exePath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("failed to get executable path: %w", err)
		}

		savePath = filepath.Dir(exePath)
	}

	err := os.MkdirAll(savePath, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to init SaveEngine: %w", err)
	}

	return &FSEngine{
		savePath: savePath,
	}, nil
}

// Saves cat image into directory of running program or into savePath.
// The supplied reader must already contain the image bytes to be written.
func (e *FSEngine) SaveImage(r io.ReadCloser, pathToFile ...string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("image reader is nil")
	}
	defer r.Close()
	if len(pathToFile) == 0 {
		return "", fmt.Errorf("missing file path")
	}

	fullpath := filepath.Join(e.savePath, filepath.Join(pathToFile...))
	if err := os.MkdirAll(filepath.Dir(fullpath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directories for %s: %w", fullpath, err)
	}

	file, err := os.Create(fullpath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", fullpath, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, r); err != nil {
		return "", fmt.Errorf("failed to save file %s: %w", fullpath, err)
	}

	return fullpath, nil
}

func (e *FSEngine) DeleteImage(pathToFile ...string) error {
	if len(pathToFile) == 0 {
		return fmt.Errorf("missing file path")
	}

	fullpath := filepath.Join(e.savePath, filepath.Join(pathToFile...))
	if err := os.Remove(fullpath); err != nil {
		return fmt.Errorf("failed to remove file %s: %w", fullpath, err)
	}

	return nil
}
