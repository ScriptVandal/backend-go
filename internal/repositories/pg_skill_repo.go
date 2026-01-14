package repositories

import (
    "database/sql"

    "github.com/ScriptVandal/backend-go/internal/models"
)

type PGSkillRepository struct {
    db *sql.DB
}

func NewPGSkillRepository(db *sql.DB) *PGSkillRepository {
    return &PGSkillRepository{db: db}
}

func (r *PGSkillRepository) List() ([]models.Skill, error) {
    rows, err := r.db.Query(`SELECT id, name, level, category FROM skills`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []models.Skill
    for rows.Next() {
        var s models.Skill
        if err := rows.Scan(&s.ID, &s.Name, &s.Level, &s.Category); err != nil {
            return nil, err
        }
        items = append(items, s)
    }
    return items, rows.Err()
}

func (r *PGSkillRepository) GetByID(id string) (*models.Skill, error) {
    var s models.Skill
    err := r.db.QueryRow(`SELECT id, name, level, category FROM skills WHERE id = $1`, id).
        Scan(&s.ID, &s.Name, &s.Level, &s.Category)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &s, nil
}

func (r *PGSkillRepository) Create(skill *models.Skill) error {
    query := `INSERT INTO skills (id, name, level, category) VALUES ($1, $2, $3, $4)`
    _, err := r.db.Exec(query, skill.ID, skill.Name, skill.Level, skill.Category)
    return err
}

func (r *PGSkillRepository) Update(skill *models.Skill) error {
    query := `UPDATE skills SET name = $2, level = $3, category = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
    _, err := r.db.Exec(query, skill.ID, skill.Name, skill.Level, skill.Category)
    return err
}

func (r *PGSkillRepository) Delete(id string) error {
    _, err := r.db.Exec(`DELETE FROM skills WHERE id = $1`, id)
    return err
}