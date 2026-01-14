package services

import "github.com/ScriptVandal/backend-go/internal/models"

type ContactRepo interface {
    List() ([]models.Contact, error)
}

type ContactService struct {
    repo ContactRepo
}

func NewContactService(repo ContactRepo) *ContactService {
    return &ContactService{repo: repo}
}

func (s *ContactService) ListContacts() ([]models.Contact, error) {
    return s.repo.List()
}
