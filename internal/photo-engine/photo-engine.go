// Package provides ready functions which will simplify action of fetching and saving images of cat via libs
// github.com/ayayaakasvin/cat-photo-fetch
package photoengine

import (
	"fmt"
	"os"
	"path/filepath"

	catphotofetch "github.com/ayayaakasvin/cat-photo-fetch"
)

// i have an idea where this worker will fetch and save images, while being exposed as main server api like */api/fetch
// where we will have directory that stores cat images accessible via web(custom viewer or using nginx)
type Worker struct {
	
}

// Saves cat image into directory of running programm or into savePath
func SaveCatImage(uuid string, savePath string) (string, error) {
	catImageBuffer, err := catphotofetch.FetchRandomPhoto()
	if err != nil {
		return "", fmt.Errorf("failed to fetch random cat image: %s", err.Error())
	}

	if savePath == "" {
		exePath, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("failed to get executable path: %w", err)
		}

		savePath = filepath.Dir(exePath)
	}

	err = os.MkdirAll(savePath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	filename := uuid + ".jpg"
	fullpath := filepath.Join(savePath, filename)

	file, err := os.Create(fullpath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %s:%s", filename, err.Error())
	}

	_, err = file.Write(catImageBuffer)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %s:%s", filename, err.Error())
	}

	return fullpath, nil
}

