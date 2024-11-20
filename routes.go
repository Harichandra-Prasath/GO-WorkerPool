package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func (S *Server) ConfigureRoutes() {
	S.Mux.Handle("/start_pool", StartPool(S))
	S.Mux.Handle("/add_job", AddJob(S))
	S.Mux.Handle("/compare", Compare())
}

func Compare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.WriteHeader(400)
			w.Write([]byte("Only GET Method is allowed\n"))
			return
		}

		Bench()
	}
}

func AddJob(S *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(400)
			w.Write([]byte("Only POST Method is allowed\n"))
			return
		}

		j := Job{}
		j.ID = uuid.New()

		raw_data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("Error in Reading the data:", err)
			w.Write([]byte("Some Error Happened\n"))
			return
		}

		err = json.Unmarshal(raw_data, &j)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("Error in Unmarshal:", err)
			w.Write([]byte("Some Error Happened\n"))
			return
		}

		go S.Pool.AddJob(&j)

	}
}

func StartPool(S *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(400)
			w.Write([]byte("Only GET Method is allowed\n"))
			return
		}
		S.Pool.Start()
	}
}
