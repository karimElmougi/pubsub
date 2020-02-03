package pubsub

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	http.ServeMux
	upgrader websocket.Upgrader
	subgroup SubscriberGroup
}

func NewServer(subgroup SubscriberGroup) *Server {
	server := &Server{
		upgrader: websocket.Upgrader{},
		subgroup: subgroup,
	}

	server.HandleFunc("/subscribe", server.handleSubscribe)
	server.HandleFunc("/publish", server.handlePublish)

	return server
}

func (h *Server) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		// websocket.Upgrader.Upgrade() already notifies the client of the error
		return
	}

	h.subgroup.AddSubscriber(conn)
}

func (h *Server) handlePublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	content := make(map[string]interface{})
	err = json.Unmarshal(body, &content)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.subgroup.Publish(content)
	w.WriteHeader(http.StatusOK)
}
