package handler

import "github.com/elastic/go-elasticsearch/v7"

type Handler struct {
	ESClient *elasticsearch.Client
}
