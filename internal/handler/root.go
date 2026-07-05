package handler

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

//go:embed templates/*.html
var templateFS embed.FS

type AllMetricsProcessor interface {
	GetAll() (map[string]map[string]string, error)
}

type Root struct {
	processor AllMetricsProcessor
	tmpl      *template.Template
}

func NewRoot(processor AllMetricsProcessor) *Root {
	return &Root{
		processor: processor,
		tmpl:      template.Must(template.ParseFS(templateFS, "templates/root.html")),
	}
}

func (v *Root) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	data, err := v.processor.GetAll()
	if err != nil {
		http.Error(res, "Load metrics data", http.StatusInternalServerError)
		return
	}

	if err := v.tmpl.Execute(res, data); err != nil {
		log.Printf("Error: %s", err.Error())
		http.Error(res, "Load data to template Error", http.StatusInternalServerError)
		return
	}
}
