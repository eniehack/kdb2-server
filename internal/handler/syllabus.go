package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/until-tsukuba/kdb2-server/internal/elasticsearch"
)

type SyllabusPayload struct {
	Id             string        `json:"id"`
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

func (h *Handler) SyllabusJSON(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "courseID")
	res, err := h.ESClient.Search(
		h.ESClient.Search.WithContext((r.Context())),
		h.ESClient.Search.WithIndex("kdb2"),
		h.ESClient.Search.WithBody(elasticsearch.NewQueryCourseIDTermQuery(query).Build()),
		h.ESClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Fatalf("ESClient err: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	response := new(ElasticSearchResponse)
	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		log.Fatalf("JSONDecoder err: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var payload SyllabusPayload
	esItem := response.Hits.Hits[0]
	payload = SyllabusPayload{
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}
