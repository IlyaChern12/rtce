package models

import (
	"sync"
	"time"
)

type Document struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Title     string    `db:"title" json:"title"`
	Body      string    `db:"body" json:"body"` // хранение текущего текста для базы
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// CRDT-поля для real-time редактирования (не сохраняются в базе)
	CRDT    []Char     `json:"crdt,omitempty"`
	mu      sync.Mutex `json:"-"`
	Version int64      `json:"-"`
}