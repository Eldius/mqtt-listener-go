package stormpersistence

import (
	"log"
	"os"
	"time"

	"github.com/Eldius/mqtt-listener-go/model"
	"github.com/Eldius/mqtt-listener-go/persistence"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
)

/*
StormRepository Repository implementation using local file storage
*/
type StormRepository struct {
	persistence.Repository
	db *storm.DB
}

/*
NewStormRepository returns a new repo
*/
func NewStormRepository() persistence.Repository {
	return &StormRepository{
		db: getDB(),
	}
}

func getDB() *storm.DB {
	_ = os.MkdirAll("./db", os.ModePerm)
	db, err := storm.Open("db/mqtt_data.db")
	if err != nil {
		log.Println("Failed to open db file")
		panic(err.Error())
	}
	err = db.Init(&model.Entry{})
	if err != nil {
		log.Println("Failed to index db")
		panic(err.Error())
	}
	return db
}

/*
Persist persists the entry
*/
func (r *StormRepository) Persist(e *model.Entry) (*model.Entry, error) {
	err := r.db.Save(e)
	if err != nil {
		log.Println("Failed to persist data to db")
		log.Println(err.Error())
	}
	return e, err
}

/*
List returns all entries for the given topic
*/
func (r *StormRepository) List(topic string) ([]*model.Entry, error) {
	results := make([]*model.Entry, 0)

	err := r.db.Find("Topic", topic, &results)
	if err != nil {
		log.Println("Failed to fetch data from db")
		log.Println(err.Error())
	}

	return results, err
}

/*
ListSince returns all entries after 'since' timestamp
*/
func (r *StormRepository) ListSince(topic string, since time.Time) ([]*model.Entry, error) {
	results := make([]*model.Entry, 0)

	constraints := q.And(
		q.Eq("Topic", topic),
		q.Gte("Timestamp", since),
	)

	query := r.db.Select(constraints)
	err := query.Find(&results)
	if err != nil {
		log.Println("Failed to fetch data from db")
		log.Println(err.Error())
	}

	return results, err
}

/*
ListLastN returns the last N results
*/
func (r *StormRepository) ListLastN(topic string, count int) ([]*model.Entry, error) {
	results := make([]*model.Entry, 0)

	constraints := q.And(
		q.Eq("Topic", topic),
	)

	err := r.db.Select(constraints).
		OrderBy("Timestamp").
		Reverse().
		Limit(count).
		Find(&results)
	if err != nil {
		log.Println("Failed to fetch data from db")
		log.Println(err.Error())
	}

	return results, err
}

/*
ListSinceNetworkEntriesGroupingByServer returns entries grouped by server name
*/
func (r *StormRepository) ListSinceNetworkEntriesGroupingByServer(topic string, since time.Time) (map[string][]*model.Entry, error) {
	var entries []*model.Entry
	results := make(map[string][]*model.Entry)

	constraints := q.And(
		q.Eq("Topic", topic),
		q.Eq("Payload.type", "network"),
		q.Gte("Timestamp", since),
	)

	err := r.db.Select(constraints).
		OrderBy("Timestamp").
		Reverse().
		Find(&entries)
	if err != nil {
		log.Println("Failed to fetch data from db")
		log.Println(err.Error())
	}

	for _, e := range entries {
		serverName := e.Payload["server"].(map[string]interface{})["name"].(string)
		results[serverName] = append(results[serverName], e)
	}
	return results, nil
}

/*
ListSinceNetworkEntriesGroupingByServer returns entries grouped by server name
*/
func (r *StormRepository) ListNetworkEntriesGroupingByServer(topic string) (map[string][]*model.Entry, error) {
	var entries []*model.Entry
	results := make(map[string][]*model.Entry)

	err := r.db.Find("Topic", topic, &entries)
	if err != nil {
		log.Println("Failed to fetch data from db")
		log.Println(err.Error())
	}

	for _, e := range entries {
		server := e.Payload["server"]
		if server != nil {
			serverName := server.(map[string]interface{})["name"].(string)
			results[serverName] = append(results[serverName], e)
		}
	}
	return results, nil
}

/*
func BackupHandleFunc(w http.ResponseWriter, _ *http.Request) {
	err := db.Bolt.View(func(tx *bbolt.Tx) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="my.db"`)
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(w)
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
*/
