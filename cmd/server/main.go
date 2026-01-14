package main

import (
    "log"
    "net/http"

    "github.com/ScriptVandal/backend-go/internal/handlers"
    "github.com/ScriptVandal/backend-go/internal/middleware"
    "github.com/ScriptVandal/backend-go/internal/repositories"
    "github.com/ScriptVandal/backend-go/internal/services"
)

func main() {
    mux := http.NewServeMux()

    // JSON repositories (file-based storage)
    projectRepo := repositories.NewJSONProjectRepository("data/projects.json")
    skillRepo := repositories.NewJSONSkillRepository("data/skills.json")
    contactRepo := repositories.NewJSONContactRepository("data/contacts.json")
    postRepo := repositories.NewJSONPostRepository("data/posts.json")

    // Services
    projectSvc := services.NewProjectService(projectRepo)
    skillSvc := services.NewSkillService(skillRepo)
    contactSvc := services.NewContactService(contactRepo)
    postSvc := services.NewPostService(postRepo)

    // Handlers
    mux.HandleFunc("/health", handlers.Health)
    mux.HandleFunc("/api/projects", handlers.NewProjectHandler(projectSvc).List)
    mux.HandleFunc("/api/skills", handlers.NewSkillHandler(skillSvc).List)
    mux.HandleFunc("/api/contacts", handlers.NewContactHandler(contactSvc).List)
    mux.HandleFunc("/api/posts", handlers.NewPostHandler(postSvc).List)

    handler := middleware.Logging(middleware.CORS(mux))

    addr := ":8080"
    log.Printf("listening on %s", addr)
    log.Fatal(http.ListenAndServe(addr, handler))
}
