package broker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"

	"github.com/kerbrek/wb-zero/app/cache"
	"github.com/kerbrek/wb-zero/app/config"
	"github.com/kerbrek/wb-zero/app/model"
	db "github.com/kerbrek/wb-zero/app/store"
)

// Wait for broker connection in a loop for the given timeout
func connectLoop(clusterID, clientID, DSN string, timeout time.Duration) (stan.Conn, error) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("connectLoop: broker connection failed after %s timeout", timeout)

		case <-ticker.C:
			sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(DSN))
			if err == nil {
				return sc, nil
			}
		}
	}
}

func MakeStan() (stan.Subscription, stan.Conn, error, error) {
	sc, connErr := connectLoop(
		config.Stan.ClusterID,
		config.Stan.ClientID,
		config.Stan.ConnStr,
		15*time.Second,
	)
	if connErr != nil {
		return nil, nil, nil, fmt.Errorf("MakeStan: %w", connErr)
	}

	sub, subErr := sc.Subscribe(
		config.Stan.Channel,
		msgHandler,
		stan.DurableName(config.Stan.DurableName),
		stan.SetManualAckMode(),
	)
	if subErr != nil {
		return nil, sc, fmt.Errorf("MakeStan: %w", subErr), nil
	}

	return sub, sc, nil, nil
}

func msgHandler(m *stan.Msg) {
	var order = new(model.Order)
	if err := json.Unmarshal(m.Data, order); err != nil {
		log.Warnf("msgHandler: %v: %s", err, m.Data)
		ackErr := m.Ack()
		if ackErr != nil {
			log.Errorf("msgHandler: %v", ackErr)
		}
		return
	}

	if err := model.Validate.Struct(order); err != nil {
		log.Warnf("msgHandler: invalid order: %s", m.Data)
		ackErr := m.Ack()
		if ackErr != nil {
			log.Errorf("msgHandler: %v", ackErr)
		}
		return
	}

	if _, ok := cache.GetOrder(order.Uid); ok {
		log.Warnf("msgHandler: duplicate order: %s", m.Data)
		ackErr := m.Ack()
		if ackErr != nil {
			log.Errorf("msgHandler: %v", ackErr)
		}
		return
	}

	compactJson, err := json.Marshal(order)
	if err != nil {
		log.Errorf("msgHandler: %v: %s", err, m.Data)
		return
	}

	if err := db.SaveOrder(order.Uid, compactJson); err != nil {
		// var pqErr *pq.Error
		// if errors.As(err, &pqErr) {
		// 	switch pqErr.Code {
		// 	case "23505": // unique_violation
		// 		// duplicate order
		// 		log.Warnf("msgHandler: %v: %s", err, m.Data)
		// 		ackErr := m.Ack()
		// 		if ackErr != nil {
		// 			log.Errorf("msgHandler: %v", ackErr)
		// 		}
		// 	default:
		// 		log.Errorf("msgHandler: %v: %s", err, m.Data)
		// 	}

		// 	return
		// }

		log.Errorf("msgHandler: %v: %s", err, m.Data)
		return
	}

	cache.SetOrder(order.Uid, compactJson)
	ackErr := m.Ack()
	if ackErr != nil {
		log.Errorf("msgHandler: %v", ackErr)
	}
}
