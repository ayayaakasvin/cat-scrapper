package domain

import (
	"context"
)

type FileMetaDataRepository interface {
	GetAllRecords(context.Context) ([]*FileMetaData, error)
	GetAllIDs(context.Context) ([]string, error)
	SaveRecord(context.Context, string, *FileMetaData) error
	GetByID(context.Context, string) (fp *FileMetaData, err error)
	DeleteRecord(string) error

	Close() error
}
