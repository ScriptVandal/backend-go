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

func (r *PGProjectRepository) GetByID(id string) (*models.Project, error) {
    var p models.Project
    var tags []string
    err := r.db.QueryRow(`SELECT id, title, description, tags, url FROM projects WHERE id = $1`, id).
        Scan(&p.ID, &p.Title, &p.Description, pq.Array(&tags), &p.URL)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    p.Tags = tags
    return &p, nil
}

func (r *PGProjectRepository) Create(project *models.Project) error {
    query := `INSERT INTO projects (id, title, description, tags, url) VALUES ($1, $2, $3, $4, $5)`
    _, err := r.db.Exec(query, project.ID, project.Title, project.Description, pq.Array(project.Tags), project.URL)
    return err
}

func (r *PGProjectRepository) Update(project *models.Project) error {
    query := `UPDATE projects SET title = $2, description = $3, tags = $4, url = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
    _, err := r.db.Exec(query, project.ID, project.Title, project.Description, pq.Array(project.Tags), project.URL)
    return err
}

func (r *PGProjectRepository) Delete(id string) error {
    _, err := r.db.Exec(`DELETE FROM projects WHERE id = $1`, id)
    return err
}
