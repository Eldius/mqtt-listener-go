package persistence

import (
	"log"
	"os"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
)

var db *storm.DB

func init() {
	_ = os.MkdirAll("./db", os.ModePerm)
	var err error
	db, err = storm.Open("db/mqtt_data.db")
	if err != nil {
		log.Println("Failed to open db file")
		panic(err.Error())
	}
	err = db.Init(&Entry{})
	if err != nil {
		log.Println("Failed to index db")
		panic(err.Error())
	}
}

func Persist(e *Entry) (*Entry, error) {
	err := db.Save(e)
	if err != nil {
		log.Println("Failed to persist data to db")
		log.Println(err.Error())
	}
	return e, err
}

func List(topic string) ([]*Entry, error) {
	results := make([]*Entry, 0)

	err := db.Find("Topic", topic, &results)
	if err != nil {
		log.Println("Failed to fetch data from db")
		log.Println(err.Error())
	}

	return results, err
}

func ListSince(topic string, since time.Time) ([]*Entry, error) {
	results := make([]*Entry, 0)

	constraints := q.And(
		q.Eq("Topic", topic),
		q.Gte("Timestamp", since),
	)

	query := db.Select(constraints)
	err := query.Find(&results)
	if err != nil {
		log.Println("Failed to fetch data from db")
		log.Println(err.Error())
	}

	return results, err
}

func ListLastN(topic string, count int) ([]*Entry, error) {
	results := make([]*Entry, 0)

	constraints := q.And(
		q.Eq("Topic", topic),
	)

	err := db.Select(constraints).
		OrderBy("Timestamp").
		Limit(count).
		Find(&results)
	if err != nil {
		log.Println("Failed to fetch data from db")
		log.Println(err.Error())
	}

	return results, err
}
