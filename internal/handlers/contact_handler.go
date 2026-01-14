package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/ScriptVandal/backend-go/internal/models"
    "github.com/ScriptVandal/backend-go/internal/repositories"
    "github.com/ScriptVandal/backend-go/internal/services"
)

type ContactHandler struct {
    svc *services.ContactService
}

func NewContactHandler(svc *services.ContactService) *ContactHandler {
    return &ContactHandler{svc: svc}
}

func (h *ContactHandler) List(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    items, err := h.svc.ListContacts()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(items)
}

func (h *ContactHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
    // Extract ID from path: /api/contacts/{id}
    path := strings.TrimPrefix(r.URL.Path, "/api/contacts/")
    id := strings.Split(path, "/")[0]

    if id == "" {
        http.Error(w, "ID is required", http.StatusBadRequest)
        return
    }

    switch r.Method {
    case http.MethodGet:
        h.Get(w, r, id)
    case http.MethodPut:
        h.Update(w, r, id)
    case http.MethodDelete:
        h.Delete(w, r, id)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func (h *ContactHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
    item, err := h.svc.GetContact(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if item == nil {
        http.Error(w, "contact not found", http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(item)
}

func (h *ContactHandler) Create(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var contact models.Contact
    if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if err := h.svc.CreateContact(&contact); err != nil {
        if err == repositories.ErrReadOnly {
            http.Error(w, "write operations not supported in JSON mode", http.StatusForbidden)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(contact)
}

func (h *ContactHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
    var contact models.Contact
    if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    contact.ID = id

    if err := h.svc.UpdateContact(&contact); err != nil {
        if err == repositories.ErrReadOnly {
            http.Error(w, "write operations not supported in JSON mode", http.StatusForbidden)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(contact)
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
    if err := h.svc.DeleteContact(id); err != nil {
        if err == repositories.ErrReadOnly {
            http.Error(w, "write operations not supported in JSON mode", http.StatusForbidden)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
