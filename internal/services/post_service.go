package services

import "github.com/ScriptVandal/backend-go/internal/models"

type PostRepo interface {
    List() ([]models.Post, error)
    GetByID(id string) (*models.Post, error)
    Create(post *models.Post) error
    Update(post *models.Post) error
    Delete(id string) error
}

type PostService struct {
    repo PostRepo
}

func NewPostService(repo PostRepo) *PostService {
    return &PostService{repo: repo}
}

func (s *PostService) ListPosts() ([]models.Post, error) {
    return s.repo.List()
}

func (s *PostService) GetPost(id string) (*models.Post, error) {
    return s.repo.GetByID(id)
}

func (s *PostService) CreatePost(post *models.Post) error {
    return s.repo.Create(post)
}

func (s *PostService) UpdatePost(post *models.Post) error {
    return s.repo.Update(post)
}

func (s *PostService) DeletePost(id string) error {
    return s.repo.Delete(id)
}