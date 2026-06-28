package domain

import (
	"io"
)

// reformat with more generic methods
type ImageFileSystem interface {
	SaveImage(io.ReadCloser, ...string) (string, error)
	DeleteImage(...string) (error)
}