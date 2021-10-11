package persistence

import (
	"time"

	"github.com/Eldius/mqtt-listener-go/model"
)

type Repository interface {
	Persist(*model.Entry) (*model.Entry, error)
	List(topic string) ([]*model.Entry, error)
	ListSince(topic string, since time.Time) ([]*model.Entry, error)
	ListLastN(topic string, count int) ([]*model.Entry, error)
	ListLastNGroupingBy(topic string, groupByField string, count int) (map[string][]*model.Entry, error)
	ListSinceNetworkEntriesGroupingByServer(topic string, since time.Time) (map[string][]*model.Entry, error)
	ListNetworkEntriesGroupingByServer(topic string) (map[string][]*model.Entry, error)
}
