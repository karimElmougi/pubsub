package pubsub

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
)

type SubscriberGroup struct {
	subChan chan *websocket.Conn
	pubChan chan map[string]interface{}
}

func NewSubscriberGroup() SubscriberGroup {
	return NewSubscriberGroupWithContext(context.Background())
}

func NewSubscriberGroupWithContext(ctx context.Context) SubscriberGroup {
	subChan := make(chan *websocket.Conn)
	pubChan := make(chan map[string]interface{})
	go subscriberGroupActor(ctx, subChan, pubChan)

	return SubscriberGroup{subChan, pubChan}
}

func (s *SubscriberGroup) Publish(content map[string]interface{}) {
	s.pubChan <- content
}

func (s *SubscriberGroup) AddSubscriber(conn *websocket.Conn) {
	s.subChan <- conn
}

func subscriberGroupActor(ctx context.Context, subChan chan *websocket.Conn, pubChan chan map[string]interface{}) {
	subscribers := make(map[*websocket.Conn]struct{})

	for {
		select {
		case <-ctx.Done():
			return
		case sub := <-subChan:
			subscribers[sub] = struct{}{}
		case content := <-pubChan:
			publish(subscribers, content)
		}
	}
}

func publish(subscribers map[*websocket.Conn]struct{}, content map[string]interface{}) {
	var toRemove []*websocket.Conn
	for s, _ := range subscribers {
		err := s.WriteJSON(content)
		if err != nil {
			toRemove = append(toRemove, s)
			log.Println("error while publishing: ", err)
		}
	}

	for _, r := range toRemove {
		delete(subscribers, r)
	}
}
