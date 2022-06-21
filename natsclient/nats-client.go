package natsclient

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	c             *nats.Conn
	subs          map[string]*nats.Subscription
	uniqueReplyTo map[string]string
	reqTimeout    time.Duration
}

// setting.json 에서 host, port 받아서 처리하기
func NewNatsClient() *NatsClient {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Println(err)
	}
	return &NatsClient{
		c:             nc,
		subs:          make(map[string]*nats.Subscription),
		uniqueReplyTo: make(map[string]string),
		reqTimeout:    time.Second * 5,
	}
}

func (n *NatsClient) Close() {
	n.c.Close()
}

func (nc *NatsClient) SetRequestTimeout(timeout time.Duration) {
	nc.reqTimeout = timeout
}

func Request[T any](nc *NatsClient, subject string) ([]T, error) {
	var uniqueReplyTo string
	var e error
	sub, ok := nc.subs[subject]
	if !ok {
		uniqueReplyTo = nats.NewInbox()
		sub, e = nc.c.SubscribeSync(uniqueReplyTo)
		if e != nil {
			log.Println(e)
			return nil, e
		}
		nc.subs[subject] = sub
		nc.uniqueReplyTo[subject] = uniqueReplyTo
	} else {
		uniqueReplyTo = nc.uniqueReplyTo[subject]
	}
	if e = nc.c.PublishRequest(subject, uniqueReplyTo, []byte("req")); e != nil {
		log.Println(e)
		return nil, e
	}
	start := time.Now()
	var responses []T
	for time.Since(start) < nc.reqTimeout {
		// Probably at maximum 0.1 seconds is enough to get the responses
		// but we'll wait for 1 second just in case it takes longer
		msg, err := sub.NextMsg(time.Second)
		if err != nil {
			e = err
			break
		}
		var response T
		if err := json.Unmarshal(msg.Data, &response); err != nil {
			e = err
			log.Println(err)
			break
		}
		responses = append(responses, response)
	}
	return responses, e
}

func (nc *NatsClient) Reply(subject string, data interface{}, wg *sync.WaitGroup) error {
	sub, ok := nc.subs[subject]
	if !ok {
		if _, err := nc.c.Subscribe(subject, func(m *nats.Msg) {
			j, _ := json.Marshal(data)
			m.Respond(j)
		}); err != nil {
			wg.Done()
			return err
		}
		nc.subs[subject] = sub
	}
	return nil
}

func Subscribe[T any](nc *NatsClient, subject string, ch chan T) error {
	sub, ok := nc.subs[subject]
	var err error
	if !ok {
		if _, err := nc.c.Subscribe(subject, func(m *nats.Msg) {
			var d T
			log.Println("Received:", string(m.Data))
			err = json.Unmarshal(m.Data, &d)
			if err != nil {
				log.Println(err)
				log.Println("Come here 1")
				return
			}
			ch <- d
		}); err != nil {
			log.Println("Come here 2")
			return err
		}
		log.Println("Come here 3")
		nc.subs[subject] = sub
	}
	return nil
}

func (nc *NatsClient) Publish(subject string, data interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	log.Println("Publishing:", string(j))
	return nc.c.Publish(subject, j)
}
