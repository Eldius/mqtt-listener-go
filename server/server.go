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
	"github.com/Eldius/mqtt-listener-go/model"
	"github.com/Eldius/mqtt-listener-go/mqttclient"
	"github.com/Eldius/mqtt-listener-go/persistence"
	"github.com/Eldius/mqtt-listener-go/persistence/mongopersistence"
	"github.com/Eldius/mqtt-listener-go/persistence/stormpersistence"
)

type ResultMode struct {
	Timestamp time.Time   `json:"timestamp,omitempty"`
	Error     string      `json:"error,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func ListLastEntrys(repo persistence.Repository) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
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

		results, err := repo.ListLastN(topic, qtt)
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

}

func QueryTopic(repo persistence.Repository) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		topic := q.Get("t")

		log.Println("---\nParameters:")
		for k, _ := range q {
			log.Printf("- %s: %s\n", k, q.Get(k))
		}
		var results []*model.Entry
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
			repo.ListSince(topic, parsedSince)
		}
		results, err := repo.List(topic)
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
}

func QueryNetworkTopic(repo persistence.Repository) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		topic := q.Get("t")

		log.Println("---\nParameters:")
		for k, _ := range q {
			log.Printf("- %s: %s\n", k, q.Get(k))
		}
		var results map[string][]*model.Entry
		if q.Get("s") != "" {
			var err error
			var parsedSince time.Time
			parsedSince, err = time.Parse(time.RFC3339, q.Get("since"))
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Header().Set("Content-Type", "application/json")
				json.NewEncoder(rw).Encode(&ResultMode{
					Error:     err.Error(),
					Timestamp: time.Now(),
				})
				return
			}
			results, err = repo.ListSinceNetworkEntriesGroupingByServer(topic, parsedSince)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Header().Set("Content-Type", "application/json")
				json.NewEncoder(rw).Encode(&ResultMode{
					Error:     err.Error(),
					Timestamp: time.Now(),
				})
				return
			}
		}
		results, err := repo.ListNetworkEntriesGroupingByServer(topic)
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
}

func Start(port int) {

	var repo persistence.Repository
	if config.UseMongoPersistence() {
		repo = mongopersistence.NewMongoRepository()
	} else {
		repo = stormpersistence.NewStormRepository()
	}

	go mqttclient.Connect(repo)

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static/build"))
	mux.Handle("/", fs)
	mux.HandleFunc("/query", QueryTopic(repo))
	mux.HandleFunc("/network/query", QueryNetworkTopic(repo))
	mux.HandleFunc("/last", ListLastEntrys(repo))
	//mux.HandleFunc("/dump", strm.BackupHandleFunc)

	host := fmt.Sprintf(":%d", port)

	log.Println(http.ListenAndServe(host, cors.CORS(mux)))
}
