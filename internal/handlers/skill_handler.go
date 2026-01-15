package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/ScriptVandal/backend-go/internal/models"
    "github.com/ScriptVandal/backend-go/internal/repositories"
    "github.com/ScriptVandal/backend-go/internal/services"
)

type SkillHandler struct {
    svc *services.SkillService
}

func NewSkillHandler(svc *services.SkillService) *SkillHandler {
    return &SkillHandler{svc: svc}
}

func (h *SkillHandler) List(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    items, err := h.svc.ListSkills()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(items)
}

func (h *SkillHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
    // Extract ID from path: /api/skills/{id}
    path := strings.TrimPrefix(r.URL.Path, "/api/skills/")
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

func (h *SkillHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
    item, err := h.svc.GetSkill(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if item == nil {
        http.Error(w, "skill not found", http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(item)
}

func (h *SkillHandler) Create(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var skill models.Skill
    if err := json.NewDecoder(r.Body).Decode(&skill); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if err := h.svc.CreateSkill(&skill); err != nil {
        if err == repositories.ErrReadOnly {
            http.Error(w, "write operations not supported in JSON mode", http.StatusForbidden)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(skill)
}

func (h *SkillHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
    var skill models.Skill
    if err := json.NewDecoder(r.Body).Decode(&skill); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    skill.ID = id

    if err := h.svc.UpdateSkill(&skill); err != nil {
        if err == repositories.ErrReadOnly {
            http.Error(w, "write operations not supported in JSON mode", http.StatusForbidden)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(skill)
}

func (h *SkillHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
    if err := h.svc.DeleteSkill(id); err != nil {
        if err == repositories.ErrReadOnly {
            http.Error(w, "write operations not supported in JSON mode", http.StatusForbidden)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
