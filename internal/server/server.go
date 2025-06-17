package server

import (
	"bookstore/internal/handlers"
	"bookstore/internal/services"
	"bookstore/internal/utils"
	"fmt"
	"log"
	"net/http"
)

const (
	DefaultPort = ":8080"
)

type Server struct {
	mux         *http.ServeMux
	bookService *services.BookService
	port        string
}

// NewServer creates a new server instance
func NewServer(port string) *Server {
	if port == "" {
		port = DefaultPort
	}

	return &Server{
		port: port,
	}
}

// Initialize configures the server with all services and routes
func (s *Server) Initialize() {
	// Initialize the books service
	s.bookService = services.NewBookService()

	// Configure the multiplexer
	s.mux = http.NewServeMux()

	// Configure the routes
	s.setupRoutes()
}

// setupRoutes configures all application routes
func (s *Server) setupRoutes() {
	// Configure the handlers
	bookHandlers := handlers.NewBookHandlers(s.bookService)

	// Routes to books
	s.mux.Handle("/books/", bookHandlers)
	s.mux.Handle("/books", bookHandlers)

	// Root route
	s.mux.HandleFunc("/", s.handleRoot)

	// Route to show all available routes
	s.mux.HandleFunc("/routes", s.handleRoutes)
}

// handleRoot handles the application's root route
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Book Management API",
		"version": "1.0.0",
		"server_info": map[string]interface{}{
			"port":   s.port,
			"status": "running",
		},
		"endpoints": map[string]string{
			"GET /":                  "API Information",
			"GET /routes":            "List all available routes",
			"GET /books":             "List all books",
			"GET /books/{id}":        "Search book by ID",
			"POST /books":            "Add new book",
			"PUT /books/{id}":        "Edit book",
			"PATCH /books/{id}/read": "Edit reading status",
			"DELETE /books/{id}":     "Remove book",
		},
	})
}

// handleRoutes returns all available routes
func (s *Server) handleRoutes(w http.ResponseWriter, r *http.Request) {
	routes := []map[string]string{
		{"method": "GET", "path": "/", "description": "API Information"},
		{"method": "GET", "path": "/routes", "description": "List all available routes"},
		{"method": "GET", "path": "/books", "description": "List all books"},
		{"method": "GET", "path": "/books/{id}", "description": "Search book by ID"},
		{"method": "POST", "path": "/books", "description": "Add new book"},
		{"method": "PUT", "path": "/books/{id}", "description": "Edit book"},
		{"method": "PATCH", "path": "/books/{id}/read", "description": "Edit reading status"},
		{"method": "DELETE", "path": "/books/{id}", "description": "Remove book"},
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"total_routes": len(routes),
		"routes":       routes,
	})
}

// Start starts the server
func (s *Server) Start() error {
	s.printStartupInfo()

	log.Printf("🚀 Server started successfully on port %s", s.port)
	return http.ListenAndServe(s.port, s.mux)
}

// printStartupInfo displays the server startup information
func (s *Server) printStartupInfo() {
	fmt.Printf("🚀 Server running on port %s\n", s.port)
	fmt.Println("📚 Book Management API")
	fmt.Println("📋 Available endpoints:")
	fmt.Println("   GET    /                 - API Information")
	fmt.Println("   GET    /routes           - List all routes")
	fmt.Println("   GET    /books            - List all books")
	fmt.Println("   GET    /books/{id}       - Search book by ID")
	fmt.Println("   POST   /books            - Add new book")
	fmt.Println("   PUT    /books/{id}       - Edit book")
	fmt.Println("   PATCH  /books/{id}/read  - Edit reading status")
	fmt.Println("   DELETE /books/{id}       - Remove book")
	fmt.Printf("🌐 Access: http://localhost%s\n", s.port)
	fmt.Printf("📋 routes: http://localhost%s/routes\n\n", s.port)
}

// GetMux returns the HTTP multiplexer (useful for testing)
func (s *Server) GetMux() *http.ServeMux {
	return s.mux
}

// GetBookService returns the book service (useful for testing)
func (s *Server) GetBookService() *services.BookService {
	return s.bookService
}
