package model

import (
	"time"
)

type Entry struct {
	ID        int        `storm:"id,increment"`
	Topic     string     `storm:"index"`
	Timestamp *time.Time `storm:"index"`
	Payload   map[string]interface{}
}

func NewEntry(topic string, payload map[string]interface{}) *Entry {
	t := time.Now()
	return &Entry{
		Topic:     topic,
		Timestamp: &t,
		Payload:   payload,
	}
}
