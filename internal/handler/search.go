package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/until-tsukuba/kdb2-server/internal/elasticsearch"
)

func (h *Handler) Result(w http.ResponseWriter, r *http.Request) {
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

	tmpl, err := template.ParseFiles("result.html.tmpl")
	if err != nil {
		log.Printf("template parse err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = tmpl.Execute(w, response.Hits.Hits); err != nil {
		log.Printf("template execute err: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
