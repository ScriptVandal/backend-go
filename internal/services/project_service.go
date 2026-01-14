package services

import "github.com/ScriptVandal/backend-go/internal/models"

type ProjectRepo interface {
    List() ([]models.Project, error)
}

type ProjectService struct {
    repo ProjectRepo
}

func NewProjectService(repo ProjectRepo) *ProjectService {
    return &ProjectService{repo: repo}
}

func (s *ProjectService) ListProjects() ([]models.Project, error) {
    return s.repo.List()
}
