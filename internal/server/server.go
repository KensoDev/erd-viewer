package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/kensodev/erd-viewer/internal/db"
	"github.com/kensodev/erd-viewer/internal/export"
	"github.com/kensodev/erd-viewer/web"
)

// Server handles HTTP requests for the ERD viewer
type Server struct {
	schemaData *db.SchemaData
	listener   net.Listener
}

// New creates a new HTTP server for serving the ERD
func New(data *db.SchemaData, listenAddr string) (*Server, error) {
	if listenAddr == "" {
		listenAddr = "127.0.0.1:0"
	}

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	return &Server{
		schemaData: data,
		listener:   ln,
	}, nil
}

// URL returns the server's URL
func (s *Server) URL() string {
	return fmt.Sprintf("http://%s", s.listener.Addr())
}

// Start begins serving HTTP requests (blocking)
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API endpoint for schema data
	mux.HandleFunc("/schema", s.handleSchema)

	// Export endpoints
	mux.HandleFunc("/export/drawio", s.handleExportDrawio)
	mux.HandleFunc("/export/plantuml", s.handleExportPlantUML)

	// Serve static assets
	mux.HandleFunc("/static/", s.handleStatic)

	// Serve the main HTML page
	mux.HandleFunc("/", s.handleIndex)

	return http.Serve(s.listener, mux)
}

func (s *Server) handleSchema(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.schemaData)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	data, err := web.Files.ReadFile("templates/index.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	// Remove /static/ prefix and read from static/
	path := r.URL.Path[1:] // Remove leading /
	data, err := web.Files.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Set content type based on extension
	contentType := "text/plain"
	if len(path) > 3 {
		ext := path[len(path)-3:]
		switch ext {
		case ".js":
			contentType = "application/javascript"
		case "css":
			contentType = "text/css"
		}
	}
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

type exportRequest struct {
	Tables []string `json:"tables"`
}

func (s *Server) handleExportDrawio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req exportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exporter := export.NewDrawioExporter()
	output, err := exporter.Export(s.schemaData, req.Tables)
	if err != nil {
		http.Error(w, fmt.Sprintf("Export failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Content-Disposition", "attachment; filename=\"schema.drawio.xml\"")
	w.Write([]byte(output))
}

func (s *Server) handleExportPlantUML(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req exportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exporter := export.NewPlantUMLExporter()
	output, err := exporter.Export(s.schemaData, req.Tables)
	if err != nil {
		http.Error(w, fmt.Sprintf("Export failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=\"schema.puml\"")
	w.Write([]byte(output))
}
