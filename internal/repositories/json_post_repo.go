package repositories

import (
    "encoding/json"
    "os"

    "github.com/ScriptVandal/backend-go/internal/models"
)

type PostRepository interface {
    List() ([]models.Post, error)
    GetByID(id string) (*models.Post, error)
    Create(post *models.Post) error
    Update(post *models.Post) error
    Delete(id string) error
}

// JSONPostRepository implements PostRepository using local JSON file.
type JSONPostRepository struct {
    path string
}

func NewJSONPostRepository(path string) *JSONPostRepository {
    return &JSONPostRepository{path: path}
}

func (r *JSONPostRepository) List() ([]models.Post, error) {
    f, err := os.Open(r.path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var items []models.Post
    if err := json.NewDecoder(f).Decode(&items); err != nil {
        return nil, err
    }
    return items, nil
}

func (r *JSONPostRepository) GetByID(id string) (*models.Post, error) {
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

func (r *JSONPostRepository) Create(post *models.Post) error {
    return ErrReadOnly
}

func (r *JSONPostRepository) Update(post *models.Post) error {
    return ErrReadOnly
}

func (r *JSONPostRepository) Delete(id string) error {
    return ErrReadOnly
}
