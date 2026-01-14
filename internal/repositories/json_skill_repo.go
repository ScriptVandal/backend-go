package repositories

import (
    "encoding/json"
    "os"

    "github.com/ScriptVandal/backend-go/internal/models"
)

type SkillRepository interface {
    List() ([]models.Skill, error)
    GetByID(id string) (*models.Skill, error)
    Create(skill *models.Skill) error
    Update(skill *models.Skill) error
    Delete(id string) error
}

type JSONSkillRepository struct {
    path string
}

func NewJSONSkillRepository(path string) *JSONSkillRepository {
    return &JSONSkillRepository{path: path}
}

func (r *JSONSkillRepository) List() ([]models.Skill, error) {
    f, err := os.Open(r.path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var items []models.Skill
    if err := json.NewDecoder(f).Decode(&items); err != nil {
        return nil, err
    }
    return items, nil
}

func (r *JSONSkillRepository) GetByID(id string) (*models.Skill, error) {
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

func (r *JSONSkillRepository) Create(skill *models.Skill) error {
    return ErrReadOnly
}

func (r *JSONSkillRepository) Update(skill *models.Skill) error {
    return ErrReadOnly
}

func (r *JSONSkillRepository) Delete(id string) error {
    return ErrReadOnly
}
