package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/ScriptVandal/backend-go/internal/models"
    "github.com/ScriptVandal/backend-go/internal/repositories"
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

func (h *PostHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
    // Extract ID from path: /api/posts/{id}
    path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
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

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
    item, err := h.svc.GetPost(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if item == nil {
        http.Error(w, "post not found", http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(item)
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var post models.Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if err := h.svc.CreatePost(&post); err != nil {
        if err == repositories.ErrReadOnly {
            http.Error(w, "write operations not supported in JSON mode", http.StatusForbidden)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
    var post models.Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    post.ID = id

    if err := h.svc.UpdatePost(&post); err != nil {
        if err == repositories.ErrReadOnly {
            http.Error(w, "write operations not supported in JSON mode", http.StatusForbidden)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
    if err := h.svc.DeletePost(id); err != nil {
        if err == repositories.ErrReadOnly {
            http.Error(w, "write operations not supported in JSON mode", http.StatusForbidden)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
