package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/ScriptVandal/backend-go/internal/services"
)

type PostHandler struct {
    svc *services.PostService
}

func NewPostHandler(svc *services.PostService) *PostHandler {
    return &PostHandler{svc: svc}
}

func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    items, err := h.svc.ListPosts()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(items)
}
