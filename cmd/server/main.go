package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"

	"github.com/ScriptVandal/backend-go/internal/config"
	"github.com/ScriptVandal/backend-go/internal/handlers"
	"github.com/ScriptVandal/backend-go/internal/middleware"
	"github.com/ScriptVandal/backend-go/internal/repositories"
	"github.com/ScriptVandal/backend-go/internal/services"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()

	usePG := cfg.DatabaseURL != ""

	// default JSON repositories
	var projectRepo repositories.ProjectRepository = repositories.NewJSONProjectRepository("data/projects.json")
	var skillRepo repositories.SkillRepository = repositories.NewJSONSkillRepository("data/skills.json")
	var contactRepo repositories.ContactRepository = repositories.NewJSONContactRepository("data/contacts.json")
	var postRepo repositories.PostRepository = repositories.NewJSONPostRepository("data/posts.json")

	var authService *services.AuthService

	// optional: switch to Postgres if DATABASE_URL is provided
	if usePG {
		db, err := sql.Open("postgres", cfg.DatabaseURL)
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

		// Auth only available with Postgres
		if cfg.JWTSecret == "" || cfg.JWTRefreshSecret == "" {
			log.Println("WARNING: JWT secrets not set. Authentication will not be available.")
		} else {
			userRepo := repositories.NewPGUserRepository(db)
			refreshTokenRepo := repositories.NewPGRefreshTokenRepository(db)
			authService = services.NewAuthService(userRepo, refreshTokenRepo, cfg)
			log.Println("Authentication enabled")
		}

		log.Println("Using PostgreSQL repositories")
	} else {
		log.Println("Using JSON repositories (read-only)")
	}

	// Services
	projectSvc := services.NewProjectService(projectRepo)
	skillSvc := services.NewSkillService(skillRepo)
	contactSvc := services.NewContactService(contactRepo)
	postSvc := services.NewPostService(postRepo)

	// Handlers
	projectHandler := handlers.NewProjectHandler(projectSvc)
	skillHandler := handlers.NewSkillHandler(skillSvc)
	contactHandler := handlers.NewContactHandler(contactSvc)
	postHandler := handlers.NewPostHandler(postSvc)

	// Health endpoint (no auth)
	mux.HandleFunc("/health", handlers.Health)

	// Auth endpoints (if available)
	if authService != nil {
		authHandler := handlers.NewAuthHandler(authService)
		mux.HandleFunc("/api/auth/register", authHandler.Register)
		mux.HandleFunc("/api/auth/login", authHandler.Login)
		mux.HandleFunc("/api/auth/refresh", authHandler.Refresh)
		mux.HandleFunc("/api/auth/logout", authHandler.Logout)
	}

	// Entity collection endpoints (GET public, POST requires auth)
	mux.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			projectHandler.Create(w, r)
		} else {
			projectHandler.List(w, r)
		}
	})
	mux.HandleFunc("/api/skills", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			skillHandler.Create(w, r)
		} else {
			skillHandler.List(w, r)
		}
	})
	mux.HandleFunc("/api/contacts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			contactHandler.Create(w, r)
		} else {
			contactHandler.List(w, r)
		}
	})
	mux.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			postHandler.Create(w, r)
		} else {
			postHandler.List(w, r)
		}
	})

	// Entity item endpoints (GET public, PUT/DELETE require auth)
	mux.HandleFunc("/api/projects/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/projects/") && r.URL.Path != "/api/projects/" {
			projectHandler.HandleItem(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/api/skills/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/skills/") && r.URL.Path != "/api/skills/" {
			skillHandler.HandleItem(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/api/contacts/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/contacts/") && r.URL.Path != "/api/contacts/" {
			contactHandler.HandleItem(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/posts/") && r.URL.Path != "/api/posts/" {
			postHandler.HandleItem(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	// Apply middleware
	var handler http.Handler = mux
	
	// Auth middleware (if auth is enabled)
	if authService != nil {
		handler = middleware.Auth(authService)(handler)
	}
	
	handler = middleware.Logging(middleware.CORS(cfg.CORSOrigins)(handler))

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}