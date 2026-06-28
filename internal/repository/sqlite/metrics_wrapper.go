package sqlite

import (
	"context"

	"github.com/ayayaakasvin/cat-scrapper/internal/domain"
	"github.com/ayayaakasvin/cat-scrapper/internal/metrics"
)

type MetricsWrapper struct {
	next domain.FileMetaDataRepository
}

func NewMetricsWrapper(next domain.FileMetaDataRepository) domain.FileMetaDataRepository {
	return &MetricsWrapper{
		next: next,
	}
}

func (m *MetricsWrapper) GetAllRecords(ctx context.Context) ([]*domain.FileMetaData, error) {
	metrics.ActiveDBConnections.Inc()
	defer metrics.ActiveDBConnections.Dec()

	records, err := m.next.GetAllRecords(ctx)
	if err == nil {
		metrics.DBQueries.Inc()
	}
	return records, err
}

func (m *MetricsWrapper) GetAllIDs(ctx context.Context) ([]string, error) {
	metrics.ActiveDBConnections.Inc()
	defer metrics.ActiveDBConnections.Dec()

	ids, err := m.next.GetAllIDs(ctx)
	if err == nil {
		metrics.DBQueries.Inc()
	}
	return ids, err
}

func (m *MetricsWrapper) SaveRecord(ctx context.Context, id string, fmd *domain.FileMetaData) error {
	metrics.ActiveDBConnections.Inc()
	defer metrics.ActiveDBConnections.Dec()

	err := m.next.SaveRecord(ctx, id, fmd)
	if err == nil {
		metrics.DBQueries.Inc()
	}
	return err
}

func (m *MetricsWrapper) GetByID(ctx context.Context, id string) (*domain.FileMetaData, error) {
	metrics.ActiveDBConnections.Inc()
	defer metrics.ActiveDBConnections.Dec()

	record, err := m.next.GetByID(ctx, id)
	if err == nil {
		metrics.DBQueries.Inc()
	}
	return record, err
}

func (m *MetricsWrapper) DeleteRecord(id string) error {
	metrics.ActiveDBConnections.Inc()
	defer metrics.ActiveDBConnections.Dec()

	err := m.next.DeleteRecord(id)
	if err == nil {
		metrics.DBQueries.Inc()
	}
	return err
}

func (m *MetricsWrapper) Close() error {

	return m.next.Close()
}
