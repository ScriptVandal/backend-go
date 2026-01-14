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