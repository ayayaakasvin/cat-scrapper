package sqlite

import (
	"database/sql"

	catphotofetch "github.com/ayayaakasvin/cat-photo-fetch"
	"github.com/ayayaakasvin/cat-scrapper/internal/domain"
)

type SqLite struct {
	conn *sql.Conn
}

// GetAllRecords implements [domain.FileMetaDataRepository].
func (s *SqLite) GetAllRecords() ([]*domain.FileMetaData, error) {
	panic("unimplemented")
}

// SaveRecord implements [domain.FileMetaDataRepository].
func (s *SqLite) SaveRecord(*domain.Job, *catphotofetch.Image) error {
	panic("unimplemented")
}

func NewSqliteConnection() domain.FileMetaDataRepository {
	return &SqLite{}
}
