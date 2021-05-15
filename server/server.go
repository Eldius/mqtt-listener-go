package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"log"

	"github.com/Eldius/cors-interceptor-go/cors"
	"github.com/Eldius/mqtt-listener-go/mqttclient"
	"github.com/Eldius/mqtt-listener-go/persistence"
)

func QueryTopic(rw http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	topic := q.Get("t")

	log.Println("---\nParameters:")
	for k, v := range q {
		log.Printf("- %s: %s\n", k, v)
	}
	var results []*persistence.Entry
	if q.Get("s") != "" {
		parsedSince, err := time.Parse(time.RFC3339, q.Get("since"))
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Header().Set("Content-Type", "text/plain")
			rw.Write([]byte(err.Error()))
		}
		persistence.ListSince(topic, parsedSince)
	}
	results, err := persistence.List(topic)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte(err.Error()))
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(&map[string]interface{}{
		"data": results,
	})
}

func Start(port int) {

	go mqttclient.Connect()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)
	mux.HandleFunc("/query", QueryTopic)

	host := fmt.Sprintf(":%d", port)

	log.Println(http.ListenAndServe(host, cors.CORS(mux)))
}
