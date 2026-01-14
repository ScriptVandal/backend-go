package services

import "github.com/ScriptVandal/backend-go/internal/models"

type ContactRepo interface {
    List() ([]models.Contact, error)
    GetByID(id string) (*models.Contact, error)
    Create(contact *models.Contact) error
    Update(contact *models.Contact) error
    Delete(id string) error
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

func (s *ContactService) GetContact(id string) (*models.Contact, error) {
    return s.repo.GetByID(id)
}

func (s *ContactService) CreateContact(contact *models.Contact) error {
    return s.repo.Create(contact)
}

func (s *ContactService) UpdateContact(contact *models.Contact) error {
    return s.repo.Update(contact)
}

func (s *ContactService) DeleteContact(id string) error {
    return s.repo.Delete(id)
}
