package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"

    _ "github.com/lib/pq"

    "github.com/ScriptVandal/backend-go/internal/handlers"
    "github.com/ScriptVandal/backend-go/internal/middleware"
    "github.com/ScriptVandal/backend-go/internal/repositories"
    "github.com/ScriptVandal/backend-go/internal/services"
)

func main() {
    mux := http.NewServeMux()

    usePG := os.Getenv("DATABASE_URL") != ""

    // default JSON repositories
    var projectRepo repositories.ProjectRepository = repositories.NewJSONProjectRepository("data/projects.json")
    var skillRepo repositories.SkillRepository = repositories.NewJSONSkillRepository("data/skills.json")
    var contactRepo repositories.ContactRepository = repositories.NewJSONContactRepository("data/contacts.json")
    var postRepo repositories.PostRepository = repositories.NewJSONPostRepository("data/posts.json")

    // optional: switch to Postgres if DATABASE_URL is provided
    if usePG {
        db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
        if err != nil {
            log.Fatal(err)
        }
        if err := db.Ping(); err != nil {
            log.Fatal(err)
        }
        projectRepo = repositories.NewPGProjectRepository(db)
        skillRepo = repositories.NewPGSkillRepository(db)
        contactRepo = repositories.NewPGContactRepository(db)
        postRepo = repositories.NewPGPostRepository(db)
        log.Println("Using PostgreSQL repositories")
    } else {
        log.Println("Using JSON repositories")
    }

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