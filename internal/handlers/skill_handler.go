package handlers

import (
    "encoding/json"
    "net/http"

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
