package course

import (
	"encoding/json"
	"fmt"
	"go_sample_api/pkg/meta"
	"log"
	"net/http"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
	}

	RequestBody struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}
)

func MakeEndpoints(log *log.Logger, s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(log, s),
	}
}

func makeCreateEndpoint(log *log.Logger, s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var body RequestBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    fmt.Sprintf("Invalid request format %s", err.Error()),
			})
			return
		}

		if body.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    "Name is required",
			})
			return
		}

		if body.StartDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    "StartDate is required",
			})
			return
		}

		if body.EndDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    "EndDate is required",
			})
			return
		}

		course, err := s.Create(
			body.Name,
			body.StartDate,
			body.EndDate,
		)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    err.Error(),
			})
		}

		json.NewEncoder(w).Encode(&Response{
			Status: 200,
			Data:   course,
		})
	}
}
