package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rltvty/go-home/timeline/schedule"

	"github.com/julienschmidt/httprouter"
)

func timeline(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jsonData, err := ioutil.ReadFile("./data.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error reading json input file")
	}
	streams, err := schedule.ParseJSON(jsonData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error parsing json input file")
	}
	sched := schedule.GetSchedule(*streams)

	jsonSched, err := json.Marshal(sched)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error marshalling schedule to json")
	}
	fmt.Fprint(w, string(jsonSched))
}

func main() {
	router := httprouter.New()
	router.GET("/timeline", timeline)

	log.Fatal(http.ListenAndServe(":4123", router))
}
