package repositories

import (
    "database/sql"

    "github.com/ScriptVandal/backend-go/internal/models"
)

type PGContactRepository struct {
    db *sql.DB
}

func NewPGContactRepository(db *sql.DB) *PGContactRepository {
    return &PGContactRepository{db: db}
}

func (r *PGContactRepository) List() ([]models.Contact, error) {
    rows, err := r.db.Query(`SELECT email, telegram, linkedin, github FROM contacts`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []models.Contact
    for rows.Next() {
        var c models.Contact
        if err := rows.Scan(&c.Email, &c.Telegram, &c.LinkedIn, &c.Github); err != nil {
            return nil, err
        }
        items = append(items, c)
    }
    return items, rows.Err()
}