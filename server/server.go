package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"log"

	"github.com/Eldius/cors-interceptor-go/cors"
	"github.com/Eldius/mqtt-listener-go/mqttclient"
	"github.com/Eldius/mqtt-listener-go/persistence"
)

func QueryTopic(rw http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	results, err := persistence.List(q.Get("t"))
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

	go mqttclient.Start()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)
	mux.HandleFunc("/query", QueryTopic)

	host := fmt.Sprintf(":%d", port)

	log.Println(http.ListenAndServe(host, cors.CORS(mux)))
}
