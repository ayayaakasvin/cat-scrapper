package domain

import "time"

type FileMetaData struct {
	ID string

	Filename string
	Filepath string

	Extension string
	MimeType  string

	Size int64

	Width  int
	Height int

	CreatedAt time.Time
}
