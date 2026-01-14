package repositories

import (
    "database/sql"

    "github.com/ScriptVandal/backend-go/internal/models"
    "github.com/lib/pq"
)

type PGPostRepository struct {
    db *sql.DB
}

func NewPGPostRepository(db *sql.DB) *PGPostRepository {
    return &PGPostRepository{db: db}
}

func (r *PGPostRepository) List() ([]models.Post, error) {
    rows, err := r.db.Query(`SELECT id, title, content, tags, published_at FROM posts`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []models.Post
    for rows.Next() {
        var p models.Post
        var tags []string
        if err := rows.Scan(&p.ID, &p.Title, &p.Content, pq.Array(&tags), &p.PublishedAt); err != nil {
            return nil, err
        }
        p.Tags = tags
        items = append(items, p)
    }
    return items, rows.Err()
}