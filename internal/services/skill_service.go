package services

import "github.com/ScriptVandal/backend-go/internal/models"

type SkillRepo interface {
    List() ([]models.Skill, error)
    GetByID(id string) (*models.Skill, error)
    Create(skill *models.Skill) error
    Update(skill *models.Skill) error
    Delete(id string) error
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

func (s *SkillService) GetSkill(id string) (*models.Skill, error) {
    return s.repo.GetByID(id)
}

func (s *SkillService) CreateSkill(skill *models.Skill) error {
    return s.repo.Create(skill)
}

func (s *SkillService) UpdateSkill(skill *models.Skill) error {
    return s.repo.Update(skill)
}

func (s *SkillService) DeleteSkill(id string) error {
    return s.repo.Delete(id)
}
