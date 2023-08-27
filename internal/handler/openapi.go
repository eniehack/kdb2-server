package handler

import (
	"io"
	"net/http"
	"os"
)

func OpenAPI(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Accept")
	switch contentType {
	case "application/json":
		file, err := os.Open("./openapi.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.Copy(w, file)
	case "application/yaml":
		file, err := os.Open("./openapi.yaml")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/yaml")
		io.Copy(w, file)
	}
	return
}
