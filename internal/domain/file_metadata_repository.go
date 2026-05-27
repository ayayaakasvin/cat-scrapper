package domain

import catphotofetch "github.com/ayayaakasvin/cat-photo-fetch"

type FileMetaDataRepository interface {
	GetAllRecords() ([]*FileMetaData, error)
	SaveRecord(*Job, *catphotofetch.Image) error
}
