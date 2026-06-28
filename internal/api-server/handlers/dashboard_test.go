package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"text/template"

	"github.com/ayayaakasvin/cat-scrapper/internal/domain"
)

type stubRepo struct {
	records []*domain.FileMetaData
	err     error
}

func (s stubRepo) GetAllRecords(context.Context,) ([]*domain.FileMetaData, error) {
	return s.records, s.err
}

func (s stubRepo) GetAllIDs(context.Context,) ([]string, error) {
	if s.err != nil {
		return nil, s.err
	}

	ids := make([]string, 0, len(s.records))
	for _, r := range s.records {
		if r != nil {
			ids = append(ids, r.UUID)
		}
	}
	return ids, nil
}

func (s stubRepo) SaveRecord(context.Context, string, *domain.FileMetaData) error {
	return nil
}

func (s stubRepo) GetByID(context.Context, string) (*domain.FileMetaData, error) {
	return nil, nil
}

func (s stubRepo) DeleteRecord(string) error {
	return nil
}

func (s stubRepo) Close() error {
	return nil
}

func TestDashboardHandlerReturnsServiceUnavailableWhenRepoIsNil(t *testing.T) {
	h := &Handlers{fmdr: nil, dashboardTmpl: template.Must(template.New("dashboard").Parse("ok"))}

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	rr := httptest.NewRecorder()

	h.DashboardHandler().ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "database repository unavailable") {
		t.Fatalf("expected error message, got %q", rr.Body.String())
	}
}

func TestDashboardHandlerReturnsInternalServerErrorWhenTemplateIsNil(t *testing.T) {
	h := &Handlers{fmdr: stubRepo{}, dashboardTmpl: nil}

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	rr := httptest.NewRecorder()

	h.DashboardHandler().ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "dashboard template unavailable") {
		t.Fatalf("expected error message, got %q", rr.Body.String())
	}
}
