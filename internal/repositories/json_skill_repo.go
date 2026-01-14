package repositories

import (
    "encoding/json"
    "os"

    "github.com/ScriptVandal/backend-go/internal/models"
)

type SkillRepository interface {
    List() ([]models.Skill, error)
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
