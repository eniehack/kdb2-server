package elasticsearch

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)

type QueryBuilder interface {
	Build() io.Reader
}

type QueryStringQueryInsideQueryString struct {
	Query  string   `json:"query"`
	Fields []string `json:"fields"`
}
type QueryStringQueryInsideQuery struct {
	QueryStringQueryInsideQueryString `json:"query_string"`
}

type QueryStringQuery struct {
	QueryStringQueryInsideQuery `json:"query"`
}

func NewQueryStringQuery(q string) *QueryStringQuery {
	return &QueryStringQuery{
		QueryStringQueryInsideQuery: QueryStringQueryInsideQuery{
			QueryStringQueryInsideQueryString: QueryStringQueryInsideQueryString{
				Query:  q,
				Fields: []string{"title", "summary", "className"},
			},
		},
	}
}

func (q *QueryStringQuery) Build() io.Reader {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(q); err != nil {
		log.Fatalf("QueryStringQuery encoding err: %v\n", err)
		return nil
	}
	return buf
}

type CourseIDTermQueryInsideCourseID struct {
	CourseID string `json:"courseID"`
}

type CourseIDTermQueryInsideTerm struct {
	Term CourseIDTermQueryInsideCourseID `json:"term"`
}

type CourseIDTermQuery struct {
	Query CourseIDTermQueryInsideTerm `json:"query"`
}

func NewQueryCourseIDTermQuery(q string) *CourseIDTermQuery {
	return &CourseIDTermQuery{
		Query: CourseIDTermQueryInsideTerm{
			Term: CourseIDTermQueryInsideCourseID{
				CourseID: q,
			},
		},
	}
}

func (q *CourseIDTermQuery) Build() io.Reader {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(q); err != nil {
		log.Fatalf("CourseIDTermQuery encoding err: %v\n", err)
		return nil
	}
	return buf
}
