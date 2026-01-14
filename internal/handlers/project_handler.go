package handlers

import (
    "encoding/json"
    "net/http"

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
