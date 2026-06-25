package fsengine

import (
	catphotofetch "github.com/ayayaakasvin/cat-photo-fetch"
	"github.com/ayayaakasvin/cat-scrapper/internal/domain"
	"github.com/ayayaakasvin/cat-scrapper/internal/metrics"
)

type MetricsWrapperFSE struct {
	next domain.ImageFileSystem
}

func NewMetricsWrapperFSE(next domain.ImageFileSystem) *MetricsWrapperFSE {
	return &MetricsWrapperFSE{
		next: next,
	}
}

func (m *MetricsWrapperFSE) SaveImage(img *catphotofetch.Image) (string, error) {
	path, err := m.next.SaveImage(img)
	if err == nil {
		metrics.TotalFilesSaved.Inc()
	}
	return path, err
}

func (m *MetricsWrapperFSE) DeleteImage(meta *domain.FileMetaData) error {
	err := m.next.DeleteImage(meta)
	if err == nil {
		metrics.TotalFilesDeleted.Inc()
	}
	return err
}
