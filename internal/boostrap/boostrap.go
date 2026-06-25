package boostrap

import (
	"github.com/ayayaakasvin/cat-scrapper/internal/domain"
	fsengine "github.com/ayayaakasvin/cat-scrapper/internal/fs_engine"
	"github.com/ayayaakasvin/cat-scrapper/internal/repository/sqlite"
)

func NewSqliteMetricsWrapper(dbPath string) (domain.FileMetaDataRepository, error) {
	db, err := sqlite.NewSqliteConnection(dbPath)
	if err != nil {
		return nil, err
	}
	return sqlite.NewMetricsWrapper(db), nil
}

func NewFSEMetricsWrapper(savePath string) (domain.ImageFileSystem, error) {
	fse, err := fsengine.NewFSE(savePath)
	if err != nil {
		return nil, err
	}
	return fsengine.NewMetricsWrapperFSE(fse), nil
}

// TODO: finish metrics implementation