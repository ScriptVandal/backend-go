package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/ScriptVandal/backend-go/internal/models"
    "github.com/ScriptVandal/backend-go/internal/repositories"
    "github.com/ScriptVandal/backend-go/internal/services"
)

type ProjectHandler struct {
    svc *services.ProjectService
}

func NewProjectHandler(svc *services.ProjectService) *ProjectHandler {
    return &ProjectHandler{svc: svc}
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    items, err := h.svc.ListProjects()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(items)
}

func (h *ProjectHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
    // Extract ID from path: /api/projects/{id}
    path := strings.TrimPrefix(r.URL.Path, "/api/projects/")
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

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
    item, err := h.svc.GetProject(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if item == nil {
        http.Error(w, "project not found", http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(item)
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var project models.Project
    if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if err := h.svc.CreateProject(&project); err != nil {
        if err == repositories.ErrReadOnly {
            http.Error(w, "write operations not supported in JSON mode", http.StatusForbidden)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
    var project models.Project
    if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    project.ID = id

    if err := h.svc.UpdateProject(&project); err != nil {
        if err == repositories.ErrReadOnly {
            http.Error(w, "write operations not supported in JSON mode", http.StatusForbidden)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
    if err := h.svc.DeleteProject(id); err != nil {
        if err == repositories.ErrReadOnly {
            http.Error(w, "write operations not supported in JSON mode", http.StatusForbidden)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
