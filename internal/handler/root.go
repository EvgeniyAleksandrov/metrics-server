package handler

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"
	"time"
)

//go:embed templates/*.html
var templateFS embed.FS

type AllMetricsProcessor interface {
	GetAll(ctx context.Context) (map[string]map[string]string, error)
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
	res.Header().Set("Content-Type", "text/html")

	getContext, cancel := context.WithTimeout(req.Context(), time.Duration(5)*time.Second)
	defer cancel()

	data, err := v.processor.GetAll(getContext)
	if err != nil {
		log.Println("ERROR:", err)
		http.Error(res, "Load metrics data", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)

	if err := v.tmpl.Execute(res, data); err != nil {
		log.Println("ERROR:", err)
		http.Error(res, "Load data to template Error", http.StatusInternalServerError)
		return
	}
}
