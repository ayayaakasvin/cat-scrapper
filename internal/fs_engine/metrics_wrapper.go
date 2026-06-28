package fsengine

import (
	"io"

	"github.com/ayayaakasvin/cat-scrapper/internal/domain"
	"github.com/ayayaakasvin/cat-scrapper/internal/metrics"
)

type MetricsWrapperFSE struct {
	next domain.ImageFileSystem
}

func NewMetricsWrapperFSE(next domain.ImageFileSystem) domain.ImageFileSystem {
	return &MetricsWrapperFSE{
		next: next,
	}
}

func (m *MetricsWrapperFSE) SaveImage(r io.ReadCloser, pathToFile ...string) (string, error) {
	path, err := m.next.SaveImage(r, pathToFile...)
	if err == nil {
		metrics.TotalFilesSaved.Inc()
	}
	return path, err
}

func (m *MetricsWrapperFSE) DeleteImage(pathToFile ...string) error {
	err := m.next.DeleteImage(pathToFile...)
	if err == nil {
		metrics.TotalFilesDeleted.Inc()
	}
	return err
}
