package services

import "github.com/ScriptVandal/backend-go/internal/models"

type SkillRepo interface {
    List() ([]models.Skill, error)
}

type SkillService struct {
    repo SkillRepo
}

func NewSkillService(repo SkillRepo) *SkillService {
    return &SkillService{repo: repo}
}

func (s *SkillService) ListSkills() ([]models.Skill, error) {
    return s.repo.List()
}
