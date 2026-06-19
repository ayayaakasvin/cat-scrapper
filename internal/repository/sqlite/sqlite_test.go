package sqlite

import (
	"context"
	"testing"

	catphotofetch "github.com/ayayaakasvin/cat-photo-fetch"
	"github.com/ayayaakasvin/cat-scrapper/internal/domain"
)

func TestBuildSaveRecordMetadataUsesActualFilePath(t *testing.T) {
	job := &domain.Job{ID: "job-1", From: "example.com", ImageUUID: "image-uuid"}
	img := &catphotofetch.Image{ContentType: "image/png", Data: []byte("abc")}

	filename, extension, mimeType, sizeBytes, err := buildSaveRecordMetadata(job, img, "/tmp/uploads/image-uuid.png")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if filename != "image-uuid.png" {
		t.Fatalf("expected filename to use the saved file name, got %q", filename)
	}
	if extension != "png" {
		t.Fatalf("expected extension png, got %q", extension)
	}
	if mimeType != "image/png" {
		t.Fatalf("expected mime type image/png, got %q", mimeType)
	}
	if sizeBytes != 3 {
		t.Fatalf("expected size to match the provided data bytes, got %d", sizeBytes)
	}
}

func TestGetAllIDsReturnsAllImageUUIDs(t *testing.T) {
	tmpFile := t.TempDir() + "/test.db"
	repo, err := NewSqliteConnection(tmpFile)
	if err != nil {
		t.Fatalf("failed to open sqlite connection: %v", err)
	}
	defer repo.Close()

	sqliteRepo, ok := repo.(*SqLite)
	if !ok {
		t.Fatal("expected sqlite repository implementation")
	}

	insert := `INSERT INTO images (id, filename, filepath, extension, mime_type, size_bytes, width, height, "from") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err := sqliteRepo.db.ExecContext(context.Background(), insert,
		"uuid-1", "a.jpg", "/images/a.jpg", "jpg", "image/jpeg", 123, 100, 100, "camera",
	); err != nil {
		t.Fatalf("failed to insert first record: %v", err)
	}
	if _, err := sqliteRepo.db.ExecContext(context.Background(), insert,
		"uuid-2", "b.png", "/images/b.png", "png", "image/png", 456, 200, 200, "scanner",
	); err != nil {
		t.Fatalf("failed to insert second record: %v", err)
	}

	ids, err := sqliteRepo.GetAllIDs()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("expected 2 ids, got %d", len(ids))
	}
	if ids[0] != "uuid-1" || ids[1] != "uuid-2" {
		t.Fatalf("expected uuids [uuid-1 uuid-2], got %v", ids)
	}
}
