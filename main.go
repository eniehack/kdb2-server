package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	ESClient *elasticsearch.Client
}

func main() {
	h := new(Handler)
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

	r.Get("/", index)
	r.Get("/result", h.result)
	r.Get("/api/v0/search", h.simplesearch)

	http.ListenAndServe(":3030", r)
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(w, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

type SubjectRef struct {
	CourseID string `json:"courseID"`
	Title    string `json:"title"`
}

type Subject struct {
	CourseID       string        `json:"courseID"`
	Title          string        `json:"title"`
	Credit         float32       `json:"credit"`
	Grade          int           `json:"grade"`
	Timetable      string        `json:"timeTable"`
	Books          []string      `json:"books"`
	ClassName      []string      `json:"className"`
	PlanPretopics  string        `json:"planPretopics"`
	Keywords       []string      `json:"keywords"`
	SeeAlsoSubject []*SubjectRef `json:"seeAlsoSubject"`
	Summary        string        `json:"summary"`
}

type Item struct {
	Index  string  `json:"_index"`
	Id     string  `json:"_id"`
	Score  float32 `json:"_score"`
	Source Subject `json:"_source"`
}

type HitsPayload struct {
	Hits []Item `json:"hits"`
}

type Response struct {
	Hits HitsPayload `json:"hits"`
}

func (h *Handler) result(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("q") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := r.URL.Query().Get("q")

	res, err := h.ESClient.Search(
		h.ESClient.Search.WithContext((r.Context())),
		h.ESClient.Search.WithIndex("kdb2"),
		h.ESClient.Search.WithBody(buildQuery(query)),
		h.ESClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Fatalf("ESClient err: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	response := new(Response)
	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		log.Fatalf("JSONDecoder err: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	/*
	   sort.SliceStable(response.Hits, func(i, j int) bool {
	           return response.Hits[j].Score < response.Hits[i].Score
	   })*/

	tmpl, err := template.ParseFiles("result.html.tmpl")
	if err != nil {
		log.Fatalf("template parse err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = tmpl.Execute(w, response.Hits.Hits); err != nil {
		log.Fatalf("template execute err: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

type SearchResponsePayload struct {
	Id             string        `json:"id"`
	Score          float32       `json:"score"`
	CourseID       string        `json:"courseID"`
	Title          string        `json:"title"`
	Credit         float32       `json:"credit"`
	Grade          int           `json:"grade"`
	Timetable      string        `json:"timeTable"`
	Books          []string      `json:"books"`
	ClassName      []string      `json:"className"`
	PlanPretopics  string        `json:"planPretopics"`
	Keywords       []string      `json:"keywords"`
	SeeAlsoSubject []*SubjectRef `json:"seeAlsoSubject"`
	Summary        string        `json:"summary"`
}

func (h *Handler) simplesearch(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("q") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := r.URL.Query().Get("q")

	res, err := h.ESClient.Search(
		h.ESClient.Search.WithContext((r.Context())),
		h.ESClient.Search.WithIndex("kdb2"),
		h.ESClient.Search.WithBody(buildQuery(query)),
		h.ESClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Fatalf("ESClient err: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	response := new(Response)
	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		log.Fatalf("JSONDecoder err: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	/*
	   sort.SliceStable(response.Hits, func(i, j int) bool {
	           return response.Hits[j].Score < response.Hits[i].Score
	   })*/

	var payload []SearchResponsePayload
	for _, esItem := range response.Hits.Hits {
		item := SearchResponsePayload{
			Score:          esItem.Score,
			Id:             esItem.Id,
			CourseID:       esItem.Source.CourseID,
			Title:          esItem.Source.Title,
			Credit:         esItem.Source.Credit,
			Grade:          esItem.Source.Grade,
			Timetable:      esItem.Source.Timetable,
			Books:          esItem.Source.Books,
			ClassName:      esItem.Source.ClassName,
			PlanPretopics:  esItem.Source.PlanPretopics,
			Keywords:       esItem.Source.Keywords,
			SeeAlsoSubject: esItem.Source.SeeAlsoSubject,
			Summary:        esItem.Source.Summary,
		}
		payload = append(payload, item)
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return
}

func buildQuery(q string) io.Reader {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"query":  q,
				"fields": []string{"title", "summary", "className"},
			},
		},
	}
	payload, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("query json marshal err: %v\n", err)
		return nil
	}
	buf := bytes.NewBuffer(payload)
	return buf
}
