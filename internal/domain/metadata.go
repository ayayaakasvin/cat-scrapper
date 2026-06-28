package domain

import "time"

type FileMetaData struct {
	UUID string

	Filename  string
	Filepath  string
	Thumbpath string

	Extension string
	MimeType  string

	Size int64

	Width  int
	Height int

	From string

	CreatedAt time.Time
}
