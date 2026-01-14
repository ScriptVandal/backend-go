package services

import "github.com/ScriptVandal/backend-go/internal/models"

type PostRepo interface {
    List() ([]models.Post, error)
}

type PostService struct {
    repo PostRepo
}

func NewPostService(repo PostRepo) *PostService {
    return &PostPostService{repo: repo}
}

func (s *PostService) ListPosts() ([]models.Post, error) {
    return s.repo.List()
}
