package imageprocessing

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"

	"github.com/disintegration/imaging"
)

const (
	DefaultWidth  = 80
	DefaultHeight = 80
)

func ImageToThumbnail(imageData io.Reader) ([]byte, error) {
	img, _, err := image.Decode(imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode img: %v", err)
	}

	thumb := imaging.Fill(img, DefaultWidth, DefaultHeight, imaging.Center, imaging.Lanczos)

	var buf bytes.Buffer

    err = jpeg.Encode(&buf, thumb, &jpeg.Options{
        Quality: 85,
    })
    if err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}