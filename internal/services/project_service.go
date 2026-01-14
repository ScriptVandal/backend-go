package services

import "github.com/ScriptVandal/backend-go/internal/models"

type ProjectRepo interface {
    List() ([]models.Project, error)
    GetByID(id string) (*models.Project, error)
    Create(project *models.Project) error
    Update(project *models.Project) error
    Delete(id string) error
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

func (s *ProjectService) GetProject(id string) (*models.Project, error) {
    return s.repo.GetByID(id)
}

func (s *ProjectService) CreateProject(project *models.Project) error {
    return s.repo.Create(project)
}

func (s *ProjectService) UpdateProject(project *models.Project) error {
    return s.repo.Update(project)
}

func (s *ProjectService) DeleteProject(id string) error {
    return s.repo.Delete(id)
}
