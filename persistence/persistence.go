package persistence

import (
	"log"

	"github.com/asdine/storm/v3"
)

var db *storm.DB

func init() {
	var err error
	db, err = storm.Open("mqtt_data.db")
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
