package user

import (
	"encoding/json"
	"go_sample_api/pkg/meta"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)
	Endpoints  struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	RequestBody struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	PutRequestBody struct {
		Firstname *string `json:"firstname"`
		Lastname  *string `json:"lastname"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
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
		Get:    makeGetEndpoint(log, s),
		GetAll: makeGetAllEndpoint(log, s),
		Update: makeUpdateEndpoint(log, s),
		Delete: makeDeleteEndpoint(log, s),
	}
}

func makeCreateEndpoint(log *log.Logger, s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var body RequestBody
		// Inject data in struct created: 'body'
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    err.Error(),
			})
			return
		}

		if body.Firstname == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    "firstname_is_required",
			})
			return
		}

		if body.Lastname == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    "lastname_is_required",
			})
			return
		}

		resp, err := s.Create(body.Firstname, body.Lastname, body.Email, body.Phone)
		if err != nil {
			log.Println("Error when trying to create user")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(resp)
	}
}

func makeDeleteEndpoint(log *log.Logger, s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DELETE /users")
		path := mux.Vars(r)
		id := path["id"]
		_, err := s.Get(id)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    err.Error(),
			})
			return
		}

		user, err := s.Delete(id)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func makeUpdateEndpoint(log *log.Logger, s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("PUT /users/:id")

		var body PutRequestBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    err.Error(),
			})
			return
		}

		if body.Firstname != nil && *body.Firstname == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    "firstname_is_required",
			})
			return
		}

		if body.Lastname != nil && *body.Lastname == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    "lastname_is_required",
			})
			return
		}

		path := mux.Vars(r)
		id := path["id"]
		_, err := s.Get(id)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    err.Error(),
			})
			return
		}

		user, err := s.Update(
			id,
			body.Firstname,
			body.Lastname,
			body.Email,
			body.Phone)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func makeGetEndpoint(log *log.Logger, s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET /users/:id")
		path := mux.Vars(r)
		id := path["id"]

		user, err := s.Get(id)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

func makeGetAllEndpoint(log *log.Logger, s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET /users")
		v := r.URL.Query()
		filters := Filters{
			Firstname: v.Get("firstname"),
			Lastname:  v.Get("lastname"),
		}

		limit, _ := strconv.Atoi(v.Get("limit"))
		page, _ := strconv.Atoi(v.Get("page"))

		count, err := s.Count(filters)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{
				Status: 500,
				Err:    err.Error(),
			})
			return
		}

		meta, err := meta.NewMeta(page, limit, count)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{
				Status: 500,
				Err:    err.Error(),
			})
			return
		}

		users, err := s.GetAll(filters, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err:    err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(&Response{
			Status: 200,
			Data:   users,
			Meta:   meta,
		})
	}
}
