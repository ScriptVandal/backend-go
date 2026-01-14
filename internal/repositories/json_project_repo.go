package repositories

import (
    "encoding/json"
    "os"

    "github.com/ScriptVandal/backend-go/internal/models"
)

type ProjectRepository interface {
    List() ([]models.Project, error)
}

type JSONProjectRepository struct {
    path string
}

func NewJSONProjectRepository(path string) *JSONProjectRepository {
    return &JSONProjectRepository{path: path}
}

func (r *JSONProjectRepository) List() ([]models.Project, error) {
    f, err := os.Open(r.path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var items []models.Project
    if err := json.NewDecoder(f).Decode(&items); err != nil {
        return nil, err
    }
    return items, nil
}
