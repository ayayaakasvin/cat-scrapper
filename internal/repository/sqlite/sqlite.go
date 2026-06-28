package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ayayaakasvin/cat-scrapper/internal/domain"
	_ "modernc.org/sqlite"
)

type SqLite struct {
	db *sql.DB
}

func NewSqliteConnection(dbPath string) (domain.FileMetaDataRepository, error) {
	if dbPath == "" {
		dbPath = "cat-scrapper.db"
	}

	dbPath, err := normalizeSqliteDBPath(dbPath)
	if err != nil {
		return nil, err
	}

	if dir := filepath.Dir(dbPath); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create sqlite directory %s: %w", dir, err)
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db %s: %w", dbPath, err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	if _, err := db.ExecContext(context.Background(), "PRAGMA busy_timeout = 5000"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set sqlite busy timeout: %w", err)
	}

	if _, err := db.ExecContext(context.Background(), "PRAGMA journal_mode = WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set sqlite journal mode: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping sqlite db %s: %w", dbPath, err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS images (
		id TEXT PRIMARY KEY,
		filename TEXT NOT NULL,
		filepath TEXT NOT NULL UNIQUE,
		thumbpath TEXT NOT NULL UNIQUE,
		extension TEXT,
		mime_type TEXT,
		size_bytes INTEGER NOT NULL,
		width INTEGER,
		height INTEGER,
		"from" TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`

	if _, err := db.ExecContext(context.Background(), createTable); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize sqlite schema: %w", err)
	}

	return &SqLite{db: db}, nil
}

func normalizeSqliteDBPath(dbPath string) (string, error) {
	if dbPath == "~" || strings.HasPrefix(dbPath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home directory: %w", err)
		}
		if dbPath == "~" {
			dbPath = home
		} else {
			dbPath = filepath.Join(home, dbPath[2:])
		}
	}

	dbPath = filepath.Clean(dbPath)
	if !filepath.IsAbs(dbPath) {
		absPath, err := filepath.Abs(dbPath)
		if err != nil {
			return "", fmt.Errorf("failed to resolve sqlite db path %s: %w", dbPath, err)
		}
		dbPath = absPath
	}

	return dbPath, nil
}

// DeleteRecord implements [domain.FileMetaDataRepository].
func (s *SqLite) DeleteRecord(id string) error {
	if s == nil || s.db == nil {
		return errors.New("sqlite connection is not initialized")
	}

	query := "DELETE FROM images WHERE id = ?"
	_, err := s.db.Exec(query, id)
	return err
}

func (s *SqLite) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

// GetByID implements [domain.FileMetaDataRepository]. Checks if such track exists, in order to prevent accessing non existing file
func (s *SqLite) GetByID(ctx context.Context, id string) (fp *domain.FileMetaData, err error) {
	if s == nil || s.db == nil {
		return nil, errors.New("sqlite connection is not initialized")
	}

	query := "SELECT id, filename, filepath, thumbpath, extension, mime_type, size_bytes, width, height, \"from\", created_at FROM images WHERE id = ? LIMIT 1"
	row := s.db.QueryRowContext(ctx, query, id)

	var (
		rowID     sql.NullString
		filename  sql.NullString
		filepath  sql.NullString
		thumbpath sql.NullString
		extension sql.NullString
		mimeType  sql.NullString
		sizeBytes sql.NullInt64
		width     sql.NullInt64
		height    sql.NullInt64
		from      sql.NullString
		createdAt sql.NullTime
	)

	if err = row.Scan(&rowID, &filename, &filepath, &thumbpath, &extension, &mimeType, &sizeBytes, &width, &height, &from, &createdAt); err != nil {
		return nil, err
	}

	return &domain.FileMetaData{
		UUID:      rowID.String,
		Filename:  filename.String,
		Filepath:  filepath.String,
		Thumbpath: thumbpath.String,
		Extension: extension.String,
		MimeType:  mimeType.String,
		Size:      sizeBytes.Int64,
		Width:     int(width.Int64),
		Height:    int(height.Int64),
		From:      from.String,
		CreatedAt: createdAt.Time,
	}, nil
}

// GetAllRecords implements [domain.FileMetaDataRepository]. Gets all records and returns slice
func (s *SqLite) GetAllRecords(ctx context.Context) ([]*domain.FileMetaData, error) {
	if s == nil || s.db == nil {
		return nil, errors.New("sqlite connection is not initialized")
	}

	query := "SELECT id, filename, filepath, extension, mime_type, size_bytes, width, height, created_at, \"from\" FROM images"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*domain.FileMetaData
	for rows.Next() {
		var (
			rowID     sql.NullString
			filename  sql.NullString
			filepath  sql.NullString
			extension sql.NullString
			mimeType  sql.NullString
			sizeBytes sql.NullInt64
			width     sql.NullInt64
			height    sql.NullInt64
			createdAt sql.NullTime
			from      sql.NullString
		)

		if err := rows.Scan(&rowID, &filename, &filepath, &extension, &mimeType, &sizeBytes, &width, &height, &createdAt, &from); err != nil {
			return nil, err
		}

		results = append(results, &domain.FileMetaData{
			UUID:      rowID.String,
			Filename:  filename.String,
			Filepath:  filepath.String,
			Extension: extension.String,
			MimeType:  mimeType.String,
			Size:      sizeBytes.Int64,
			Width:     int(width.Int64),
			Height:    int(height.Int64),
			From:      from.String,
			CreatedAt: createdAt.Time,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *SqLite) GetAllIDs(ctx context.Context) ([]string, error) {
	if s == nil || s.db == nil {
		return nil, errors.New("sqlite connection is not initialized")
	}

	query := "SELECT id FROM images"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id sql.NullString
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		if id.Valid {
			ids = append(ids, id.String)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

// SaveRecord implements [domain.FileMetaDataRepository]. Saves record
func (s *SqLite) SaveRecord(ctx context.Context, from string, fmd *domain.FileMetaData) error {
	if s == nil || s.db == nil {
		return errors.New("sqlite connection is not initialized")
	}
	if fmd == nil {
		return errors.New("image is nil")
	}

	query := "INSERT INTO images (id, filename, filepath, thumbpath, extension, mime_type, size_bytes, width, height, `from`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := s.db.ExecContext(ctx, query,
		fmd.UUID,
		fmd.Filename,
		fmd.Filepath,
		fmd.Thumbpath,
		fmd.Extension,
		fmd.MimeType,
		fmd.Size,
		fmd.Width,
		fmd.Height,
		from,
	)

	return err
}
