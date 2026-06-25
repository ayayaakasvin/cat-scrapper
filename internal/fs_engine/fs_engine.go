// Package provides ready functions which will simplify action of fetching and saving images of cat via libs
// github.com/ayayaakasvin/cat-photo-fetch
package fsengine

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	catphotofetch "github.com/ayayaakasvin/cat-photo-fetch"
	"github.com/ayayaakasvin/cat-scrapper/internal/domain"
)

type FSEngine struct {
	savePath  string
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
		savePath:  savePath,
	}, nil
}

// Saves cat image into directory of running program or into savePath.
// The supplied reader must already contain the image bytes to be written.
func (e *FSEngine) SaveImage(img *catphotofetch.Image) (string, error) {
	r := img.Reader()

	if r == nil {
		return "", fmt.Errorf("image reader is nil")
	}
	defer r.Close()

	ctype := strings.TrimPrefix(img.ContentType, "image/")
	filename := img.UUID + "." + ctype
	fullpath := filepath.Join(e.savePath, filename)

	file, err := os.Create(fullpath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, r); err != nil {
		return "", fmt.Errorf("failed to save file %s: %w", filename, err)
	}

	return fullpath, nil
}

func (e *FSEngine) DeleteImage(fmd *domain.FileMetaData) error {
	if err := os.Remove(fmd.Filepath); err != nil {
		return err
	}

	return nil
}