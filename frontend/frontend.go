package frontend

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/*
var FS embed.FS

// GetIndexTemplate returns the parsed index template
func GetIndexTemplate() (*template.Template, error) {
	return template.ParseFS(FS, "templates/index.html")
}

// ServeStatic handles static file serving
func ServeStatic() http.Handler {
	return http.FileServer(http.FS(FS))
}