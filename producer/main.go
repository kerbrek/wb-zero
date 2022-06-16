// Publish several messages to NATS Streaming server
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nats-io/stan.go"
)

const orderTemplate = `{
	"order_uid": "%s",
	"track_number": "WBILMTESTTRACK",
	"entry": "WBIL",
	"delivery": {
		"name": "Test Testov",
		"phone": "+9720000000",
		"zip": "2639809",
		"city": "Kiryat Mozkin",
		"address": "Ploshad Mira 15",
		"region": "Kraiot",
		"email": "test@gmail.com"
	},
	"payment": {
		"transaction": "%s",
		"request_id": "",
		"currency": "USD",
		"provider": "wbpay",
		"amount": 1817,
		"payment_dt": 1637907727,
		"bank": "alpha",
		"delivery_cost": 1500,
		"goods_total": 317,
		"custom_fee": 0
	},
	"items": [
		{
			"chrt_id": 9934930,
			"track_number": "WBILMTESTTRACK",
			"price": 453,
			"rid": "ab4219087a764ae0btest",
			"name": "Mascaras",
			"sale": 30,
			"size": "0",
			"total_price": 317,
			"nm_id": 2389212,
			"brand": "Vivienne Sabo",
			"status": 202
		}
	],
	"locale": "en",
	"internal_signature": "",
	"customer_id": "test",
	"delivery_service": "meest",
	"shardkey": "9",
	"sm_id": 99,
	"date_created": "2021-11-26T06:22:19Z",
	"oof_shard": "1"
}`

const orderIdMaxLen = 50

var invalidJson = `{`
var invalidOrder = `{"order_uid": "invalid"}`
var tooLongId = strings.Repeat("a", orderIdMaxLen+1)
var orderWithTooLongId = fmt.Sprintf(orderTemplate, tooLongId, tooLongId)

func main() {
	var n int64
	flag.Int64Var(&n, "n", 5, "number of messages")
	flag.Parse()

	clusterID := "test-cluster"
	clientID := "orders-producer-1"
	connStr := "nats://127.0.0.1:4222"
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(connStr))
	if err != nil {
		log.Fatal(err)
	}

	ns := time.Now().UnixNano()
	var i int64 = 0
	for ; i < n; i++ {
		id := fmt.Sprint(ns + i)
		err := sc.Publish("orders", []byte(fmt.Sprintf(orderTemplate, id, id)))
		if err != nil {
			log.Println(err)
		}
	}

	constId := "b563feb7b2b84b6test"
	err = sc.Publish("orders", []byte(fmt.Sprintf(orderTemplate, constId, constId)))
	if err != nil {
		log.Println(err)
	}

	err = sc.Publish("orders", []byte(invalidJson))
	if err != nil {
		log.Println(err)
	}

	err = sc.Publish("orders", []byte(invalidOrder))
	if err != nil {
		log.Println(err)
	}

	err = sc.Publish("orders", []byte(orderWithTooLongId))
	if err != nil {
		log.Println(err)
	}
}
