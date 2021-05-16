package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"log"

	"github.com/Eldius/cors-interceptor-go/cors"
	"github.com/Eldius/mqtt-listener-go/config"
	"github.com/Eldius/mqtt-listener-go/mqttclient"
	"github.com/Eldius/mqtt-listener-go/persistence"
)

type ResultMode struct {
	Timestamp time.Time   `json:"timestamp,omitempty"`
	Error     string      `json:"error,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func ListLastEntrys(rw http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	topic := q.Get("t")
	_qtt := q.Get("c")
	var qtt int
	if _qtt == "" {
		qtt = config.GetDefaultFetchCount()
	} else {
		var err error
		qtt, err = strconv.Atoi(_qtt)
		if err != nil {
			log.Printf("Failed to parse quantity (%s)\n%s\n", _qtt, err.Error())
			qtt = config.GetDefaultFetchCount()
		}
	}

	results, err := persistence.ListLastN(topic, qtt)
	if err != nil {
		log.Printf("Failed to list last ebtrys:\n%s\n", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(&ResultMode{
			Error:     err.Error(),
			Timestamp: time.Now(),
		})
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(&ResultMode{
		Timestamp: time.Now(),
		Data:      results,
	})

}

func QueryTopic(rw http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	topic := q.Get("t")

	log.Println("---\nParameters:")
	for k, _ := range q {
		log.Printf("- %s: %s\n", k, q.Get(k))
	}
	var results []*persistence.Entry
	if q.Get("s") != "" {
		parsedSince, err := time.Parse(time.RFC3339, q.Get("since"))
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Header().Set("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(&ResultMode{
				Error:     err.Error(),
				Timestamp: time.Now(),
			})
			return
		}
		persistence.ListSince(topic, parsedSince)
	}
	results, err := persistence.List(topic)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(&ResultMode{
			Error:     err.Error(),
			Timestamp: time.Now(),
		})
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(&ResultMode{
		Timestamp: time.Now(),
		Data:      results,
	})
}

func Start(port int) {

	go mqttclient.Connect()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)
	mux.HandleFunc("/query", QueryTopic)
	mux.HandleFunc("/last", ListLastEntrys)

	host := fmt.Sprintf(":%d", port)

	log.Println(http.ListenAndServe(host, cors.CORS(mux)))
}
