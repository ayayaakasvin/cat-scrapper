package domain

import "time"

type FileMetaData struct {
	UUID string

	Filename string
	Filepath string

	Extension string
	MimeType  string

	Size int64

	Width  int
	Height int

	CreatedAt time.Time
}
