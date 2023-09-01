package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/until-tsukuba/kdb2-server/internal/elasticsearch"
)

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

type ElasticSearchResponse struct {
	Hits struct {
		Hits []Item `json:"hits"`
	} `json:"hits"`
}

func (h *Handler) SimpleSearch(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("q") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := r.URL.Query().Get("q")

	res, err := h.ESClient.Search(
		h.ESClient.Search.WithContext((r.Context())),
		h.ESClient.Search.WithIndex("kdb2"),
		h.ESClient.Search.WithBody(elasticsearch.NewQueryStringQuery(query).Build()),
		h.ESClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Printf("ESClient err: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	response := new(ElasticSearchResponse)
	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		log.Printf("JSONDecoder err: %v\n", err)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}
