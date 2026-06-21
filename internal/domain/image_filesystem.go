package domain

import (
	catphotofetch "github.com/ayayaakasvin/cat-photo-fetch"
)

type ImageFileSystem interface {
	SaveImage(*catphotofetch.Image) (string, error)
	DeleteImage(*FileMetaData) (error)
}
