package repositories

import (
    "encoding/json"
    "os"

    "github.com/ScriptVandal/backend-go/internal/models"
)

type ContactRepository interface {
    List() ([]models.Contact, error)
    GetByID(id string) (*models.Contact, error)
    Create(contact *models.Contact) error
    Update(contact *models.Contact) error
    Delete(id string) error
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

func (r *JSONContactRepository) GetByID(id string) (*models.Contact, error) {
    items, err := r.List()
    if err != nil {
        return nil, err
    }
    for _, item := range items {
        if item.ID == id {
            return &item, nil
        }
    }
    return nil, nil
}

func (r *JSONContactRepository) Create(contact *models.Contact) error {
    return ErrReadOnly
}

func (r *JSONContactRepository) Update(contact *models.Contact) error {
    return ErrReadOnly
}

func (r *JSONContactRepository) Delete(id string) error {
    return ErrReadOnly
}
