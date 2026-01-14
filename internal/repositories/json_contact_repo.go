package repositories

import (
    "encoding/json"
    "os"

    "github.com/ScriptVandal/backend-go/internal/models"
)

type ContactRepository interface {
    List() ([]models.Contact, error)
}

type JSONContactRepository struct {
    path string
}

func NewJSONContactRepository(path string) *JSONContactRepository {
    return &JSONContactRepository{path: path}
}

func (r *JSONContactRepository) List() ([]models.Contact, error) {
    f, err := os.Open(r.path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var items []models.Contact
    if err := json.NewDecoder(f).Decode(&items); err != nil {
        return nil, err
    }
    return items, nil
}
