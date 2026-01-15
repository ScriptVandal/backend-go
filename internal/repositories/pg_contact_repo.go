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
    rows, err := r.db.Query(`SELECT id, email, telegram, linkedin, github FROM contacts`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []models.Contact
    for rows.Next() {
        var c models.Contact
        if err := rows.Scan(&c.ID, &c.Email, &c.Telegram, &c.LinkedIn, &c.Github); err != nil {
            return nil, err
        }
        items = append(items, c)
    }
    return items, rows.Err()
}

func (r *PGContactRepository) GetByID(id string) (*models.Contact, error) {
    var c models.Contact
    err := r.db.QueryRow(`SELECT id, email, telegram, linkedin, github FROM contacts WHERE id = $1`, id).
        Scan(&c.ID, &c.Email, &c.Telegram, &c.LinkedIn, &c.Github)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &c, nil
}

func (r *PGContactRepository) Create(contact *models.Contact) error {
    query := `INSERT INTO contacts (email, telegram, linkedin, github) VALUES ($1, $2, $3, $4)`
    _, err := r.db.Exec(query, contact.Email, contact.Telegram, contact.LinkedIn, contact.Github)
    return err
}

func (r *PGContactRepository) Update(contact *models.Contact) error {
    query := `UPDATE contacts SET email = $1, telegram = $2, linkedin = $3, github = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`
    _, err := r.db.Exec(query, contact.Email, contact.Telegram, contact.LinkedIn, contact.Github, contact.ID)
    return err
}

func (r *PGContactRepository) Delete(id string) error {
    _, err := r.db.Exec(`DELETE FROM contacts WHERE id = $1`, id)
    return err
}