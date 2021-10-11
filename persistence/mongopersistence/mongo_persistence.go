package mongopersistence

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Eldius/mqtt-listener-go/config"
	"github.com/Eldius/mqtt-listener-go/model"
	"github.com/Eldius/mqtt-listener-go/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	eventDatabaseName = "events"
)

type MongoRepository struct {
	persistence.Repository
	client *mongo.Client
}

func getClient() *mongo.Client {
	log.Println("mongo url:", config.GetMongoURL())
	client, err := mongo.NewClient(options.Client().ApplyURI(config.GetMongoURL()))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	/*
	   List databases
	*/
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	return client
}

func NewMongoRepository() persistence.Repository {
	return &MongoRepository{
		client: getClient(),
	}
}

func (r *MongoRepository) getEventsDatabase() *mongo.Database {
	return r.client.Database(eventDatabaseName)
}

func (r *MongoRepository) Persist(e *model.Entry) (*model.Entry, error) {
	s, err := r.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer s.EndSession(context.Background())

	col := r.getEventsDatabase().Collection(e.Topic)
	res, err := col.InsertOne(context.Background(), e)
	if err != nil {
		return nil, err
	}
	log.Printf("result: %s", res.InsertedID)
	return e, nil
}

func (r *MongoRepository) List(topic string) ([]*model.Entry, error) {
	var results []*model.Entry

	col := r.getEventsDatabase().Collection(topic)
	cur, err := col.Find(context.Background(), bson.D{})
	if err != nil {
		return results, err
	}
	if err := cur.All(context.Background(), &results); err != nil {
		return make([]*model.Entry, 0), err
	}

	return results, nil
}

func (r *MongoRepository) ListSince(topic string, since time.Time) ([]*model.Entry, error) {
	var results []*model.Entry

	col := r.getEventsDatabase().Collection(topic)
	//filter := bson.D{{"timestamp", bson.D{{"$gte", since}}}}
	filter := bson.D{{"timestamp", bson.D{{"$gte", since}}}}

	cur, err := col.Find(context.Background(), filter)
	if err != nil {
		return results, err
	}
	if err := cur.All(context.Background(), &results); err != nil {
		return make([]*model.Entry, 0), err
	}

	return results, nil
}

func (r *MongoRepository) ListLastN(topic string, count int) ([]*model.Entry, error) {
	return make([]*model.Entry, 0), nil
}

func (r *MongoRepository) ListSinceNetworkEntriesGroupingByServer(topic string, since time.Time) (map[string][]*model.Entry, error) {
	results := make(map[string][]*model.Entry)
	col := r.getEventsDatabase().Collection(topic)
	filter := bson.D{{"timestamp", bson.D{{"$gte", since}}}}

	cur, err := col.Find(context.Background(), filter)
	if err != nil {
		return results, err
	}

	var entries []*model.Entry
	cur.All(context.Background(), &entries)
	for _, e := range entries {
		serverName := e.Payload["server"].(map[string]interface{})["name"].(string)
		results[serverName] = append(results[serverName], e)
	}

	return results, nil
}

func (r *MongoRepository) ListNetworkEntriesGroupingByServer(topic string) (map[string][]*model.Entry, error) {
	results := make(map[string][]*model.Entry)
	col := r.getEventsDatabase().Collection(topic)

	cur, err := col.Find(context.Background(), bson.D{})
	if err != nil {
		return results, err
	}

	var entries []*model.Entry
	cur.All(context.Background(), &entries)
	for _, e := range entries {
		serverName := e.Payload["server"].(map[string]interface{})["name"].(string)
		results[serverName] = append(results[serverName], e)
	}

	return results, nil
}
