package repositories

import (
    "encoding/json"
    "os"

    "github.com/ScriptVandal/backend-go/internal/models"
)

type ProjectRepository interface {
    List() ([]models.Project, error)
    GetByID(id string) (*models.Project, error)
    Create(project *models.Project) error
    Update(project *models.Project) error
    Delete(id string) error
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

func (r *JSONProjectRepository) GetByID(id string) (*models.Project, error) {
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

func (r *JSONProjectRepository) Create(project *models.Project) error {
    return ErrReadOnly
}

func (r *JSONProjectRepository) Update(project *models.Project) error {
    return ErrReadOnly
}

func (r *JSONProjectRepository) Delete(id string) error {
    return ErrReadOnly
}
