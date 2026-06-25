package domain

import (
	"context"

	catphotofetch "github.com/ayayaakasvin/cat-photo-fetch"
)

type FileMetaDataRepository interface {
	GetAllRecords(context.Context) ([]*FileMetaData, error)
	GetAllIDs(context.Context) ([]string, error)
	SaveRecord(context.Context, string, *catphotofetch.Image, string) error
	GetByID(context.Context, string) (fp *FileMetaData, err error)
	DeleteRecord(string) error

	Close() error
}