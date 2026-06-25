package domain

import (
	catphotofetch "github.com/ayayaakasvin/cat-photo-fetch"
)

type ImageFileSystem interface {
	SaveImage(*catphotofetch.Image) (string, error)
	DeleteImage(*FileMetaData) (error)
}

type FSEStats struct {
	OverallQueries uint64
	ActiveQueries  int64
	FailedQueries  uint64
	AvgLatency     int64
}