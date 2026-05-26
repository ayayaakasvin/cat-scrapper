// Package provides ready functions which will simplify action of fetching and saving images of cat via libs
// github.com/ayayaakasvin/cat-photo-fetch
package saveengine

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	catphotofetch "github.com/ayayaakasvin/cat-photo-fetch"
)

type SaveEngine struct {
	savePath string
	sem      chan struct{}
}

func NewSaveEngine(savePath string) (*SaveEngine, error) {
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

	return &SaveEngine{
		savePath: savePath,
		sem:      make(chan struct{}, 8),
	}, nil
}

// Saves cat image into directory of running program or into savePath.
// The supplied reader must already contain the image bytes to be written.
func (e *SaveEngine) SaveCatImage(uid *Job, img *catphotofetch.Image) (string, error) {
	e.sem <- struct{}{}
	defer func() { <-e.sem }()
	r := img.Reader()

	if r == nil {
		return "", fmt.Errorf("image reader is nil")
	}
	
	ctype := strings.TrimPrefix(img.ContentType, "image/")
	filename := uid.ImageUUID + "." + ctype
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
