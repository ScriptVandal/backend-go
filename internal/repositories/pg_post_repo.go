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

func (r *PGPostRepository) GetByID(id string) (*models.Post, error) {
    var p models.Post
    var tags []string
    err := r.db.QueryRow(`SELECT id, title, content, tags, published_at FROM posts WHERE id = $1`, id).
        Scan(&p.ID, &p.Title, &p.Content, pq.Array(&tags), &p.PublishedAt)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    p.Tags = tags
    return &p, nil
}

func (r *PGPostRepository) Create(post *models.Post) error {
    query := `INSERT INTO posts (id, title, content, tags, published_at) VALUES ($1, $2, $3, $4, $5)`
    _, err := r.db.Exec(query, post.ID, post.Title, post.Content, pq.Array(post.Tags), post.PublishedAt)
    return err
}

func (r *PGPostRepository) Update(post *models.Post) error {
    query := `UPDATE posts SET title = $2, content = $3, tags = $4, published_at = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
    _, err := r.db.Exec(query, post.ID, post.Title, post.Content, pq.Array(post.Tags), post.PublishedAt)
    return err
}

func (r *PGPostRepository) Delete(id string) error {
    _, err := r.db.Exec(`DELETE FROM posts WHERE id = $1`, id)
    return err
}