package handlers

import (
    "encoding/json"
    "net/http"

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
