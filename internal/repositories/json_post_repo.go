package repositories

import (
    "encoding/json"
    "os"

    "github.com/ScriptVandal/backend-go/internal/models"
)

type PostRepository interface {
    List() ([]models.Post, error)
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
