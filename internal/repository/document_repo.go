package repository

import (
	"context"
	"database/sql"

	"github.com/IlyaChern12/rtce/internal/models"
)

// репо для документов
type DocumentRepository struct {
	db *sql.DB
}

// конструктор репо
func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{
		db: db,
	}
}

// создание документа
func (dr *DocumentRepository) Create(ctx context.Context, doc *models.Document) error {
	query := `INSERT INTO documents (id, user_id, title, body, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, NOW(), NOW())`

	_, err := dr.db.ExecContext(ctx, query, doc.ID, doc.UserID, doc.Title, doc.Body)
	return err
}

// достаем док по id
func (dr *DocumentRepository) GetByID(ctx context.Context, userID string) ([]*models.Document, error) {
	query := `SELECT id, user_id, title, body, created_at, updated_at FROM documents WHERE user_id = $1`

	rows, err := dr.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []*models.Document

	for rows.Next() {
		doc := &models.Document{}
		if err := rows.Scan(&doc.ID, &doc.UserID, &doc.Title, &doc.Body, &doc.CreatedAt, &doc.UpdatedAt); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, nil
}