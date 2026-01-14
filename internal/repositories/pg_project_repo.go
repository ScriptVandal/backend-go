package repositories

import (
    "database/sql"

    "github.com/ScriptVandal/backend-go/internal/models"
    "github.com/lib/pq"
)

type PGProjectRepository struct {
    db *sql.DB
}

func NewPGProjectRepository(db *sql.DB) *PGProjectRepository {
    return &PGProjectRepository{db: db}
}

func (r *PGProjectRepository) List() ([]models.Project, error) {
    rows, err := r.db.Query(`SELECT id, title, description, tags, url FROM projects`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []models.Project
    for rows.Next() {
        var p models.Project
        var tags []string
        if err := rows.Scan(&p.ID, &p.Title, &p.Description, pq.Array(&tags), &p.URL); err != nil {
            return nil, err
        }
        p.Tags = tags
        items = append(items, p)
    }
    return items, rows.Err()
}
