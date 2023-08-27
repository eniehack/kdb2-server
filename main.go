package main

import (
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/until-tsukuba/kdb2-server/internal/handler"
)

func main() {
	h := new(handler.Handler)
	cfg := &elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	esclient, err := elasticsearch.NewClient(*cfg)
	if err != nil {
		log.Fatalf("esclient init: %v\n", err)
		return
	}
	h.ESClient = esclient

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Get("/", handler.Index)
	r.Get("/result", h.Result)
	r.Get("/api/v0/syllabus/{courseID}.json", h.SyllabusJSON)
	r.Get("/api/v0/search", h.SimpleSearch)

	http.ListenAndServe(":3030", r)
}
