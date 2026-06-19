package domain

import catphotofetch "github.com/ayayaakasvin/cat-photo-fetch"

type FileMetaDataRepository interface {
	GetAllRecords() ([]*FileMetaData, error)
	GetAllIDs() ([]string, error)
	SaveRecord(*Job, *catphotofetch.Image, string) error
	GetByID(id string) (fp *FileMetaData, err error)

	Close() error
}
