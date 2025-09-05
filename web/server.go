package web

import (
	"net/http"
	"os"
	"path/filepath"

	"palworld-helper/internal/adapters/web/handlers"
	"palworld-helper/internal/core/ports"
)

type Server struct {
	craftingService ports.CraftingService
	adminService    ports.AdminService
}

func NewServer(craftingService ports.CraftingService, adminService ports.AdminService) *Server {
	return &Server{
		craftingService: craftingService,
		adminService:    adminService,
	}
}

func (s *Server) Start(addr string) error {
	// Initialize handlers
	craftingHandler := handlers.NewCraftingHandler(s.craftingService)
	adminHandler := handlers.NewAdminHandler(s.adminService)

	// Setup routes
	mux := http.NewServeMux()

	// Static files - chemin relatif au binaire
	staticDir := "./web/static/"

	// Debug: afficher le répertoire de travail
	wd, _ := os.Getwd()
	println("Working directory:", wd)
	println("Static directory:", filepath.Join(wd, staticDir))

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	// Main crafting interface
	mux.HandleFunc("/", craftingHandler.HomePage)
	mux.HandleFunc("/api/recipes", craftingHandler.GetRecipes)
	mux.HandleFunc("/api/calculate", craftingHandler.CalculateResources)

	// Admin interface
	mux.HandleFunc("/admin", adminHandler.AdminPage)
	mux.HandleFunc("/admin/api/schema", adminHandler.GetSchema)
	mux.HandleFunc("/admin/api/table/", adminHandler.HandleTableOperations)
	mux.HandleFunc("/admin/api/query", adminHandler.ExecuteQuery)
	mux.HandleFunc("/admin/api/create-table", adminHandler.CreateTable)

	return http.ListenAndServe(addr, mux)
}
