package pubsub_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/karimElmougi/pubsub/internal/pubsub"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestContent struct {
	Foo string `json:"foo"`
	N   int    `json:"n"`
}

var _ = Describe("pubsub", func() {
	var ctx context.Context
	var cancel context.CancelFunc
	var server *httptest.Server

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		subgroup := pubsub.NewSubscriberGroupWithContext(ctx)
		handler := pubsub.NewServer(subgroup)
		server = httptest.NewServer(handler)
	})

	AfterEach(func() {
		server.Close()
		cancel()
	})

	Describe("Subscribing", func() {
		Context("with incorrect protocol", func() {
			It("should return error", func() {
				response, err := http.Post(server.URL+"/subscribe", "application/json", nil)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).ToNot(Equal(http.StatusOK))
			})
		})

		Context("with ws protocol", func() {
			It("should be successful", func() {
				url := "ws" + strings.TrimPrefix(server.URL, "http") + "/subscribe"
				ws, _, err := websocket.DefaultDialer.Dial(url, nil)
				Expect(err).ToNot(HaveOccurred())
				err = ws.Close()
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Publishing", func() {
		Context("with subscribers", func() {
			It("should send content to all", func() {
				url := "ws" + strings.TrimPrefix(server.URL, "http") + "/subscribe"

				var subscribers []*websocket.Conn
				for i := 0; i < 5; i++ {
					ws, _, err := websocket.DefaultDialer.Dial(url, nil)
					Expect(err).ToNot(HaveOccurred())
					subscribers = append(subscribers, ws)
				}

				content := TestContent{Foo: "bar", N: 1}
				body, err := json.Marshal(content)
				Expect(err).ToNot(HaveOccurred())

				response, err := http.Post(server.URL+"/publish", "application/json", bytes.NewReader(body))
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				for _, sub := range subscribers {
					message := TestContent{}
					err = sub.ReadJSON(&message)
					Expect(err).ToNot(HaveOccurred())
					Expect(message).To(Equal(content))
				}
			})
		})
	})
})
