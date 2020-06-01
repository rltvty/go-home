package main

import (
	"github.com/gorilla/mux"
	"github.com/rltvty/go-home/logwrapper"
	"go.uber.org/zap"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	log := logwrapper.GetInstance()
	r := mux.NewRouter()
	get := r.Methods("GET").Subrouter()
	//post := r.Methods("POST").Subrouter()

	get.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello World!")
	})

	get.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("./speaker_endpoints.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(content))
	})


	log.Fatal("Error starting API server", zap.Any("error", http.ListenAndServe(":8000", r)))
}

