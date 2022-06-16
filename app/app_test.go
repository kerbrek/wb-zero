package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kerbrek/wb-zero/app/api"
	"github.com/kerbrek/wb-zero/app/broker"
	"github.com/kerbrek/wb-zero/app/cache"
	"github.com/kerbrek/wb-zero/app/config"
	"github.com/kerbrek/wb-zero/app/store"
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

var invalidOrderId = "invalid"
var invalidOrder = `{"order_uid": "invalid"}`

var tooLongId = strings.Repeat("a", orderIdMaxLen+1)
var orderWithTooLongId = fmt.Sprintf(orderTemplate, tooLongId, tooLongId)

var validOrderId = "b563feb7b2b84b6test"
var validOrder = fmt.Sprintf(orderTemplate, validOrderId, validOrderId)

func TestApiRoutes(t *testing.T) {
	assert := assert.New(t)
	_, err := store.MakeDB(config.Postgres.ConnStr)
	if err != nil {
		t.Fatal(err)
	}

	// db is empty for now
	empty, err := store.ReadAllOrders()
	if err != nil {
		t.Fatal(err)
	}
	// init empty cache
	cache.Init(empty)

	_, sc, subErr, connErr := broker.MakeStan()
	if connErr != nil {
		t.Fatal(connErr)
	}
	if subErr != nil {
		t.Fatal(subErr)
	}

	// publish invalid json
	err = sc.Publish(config.Stan.Channel, []byte(invalidJson))
	if err != nil {
		t.Fatal(err)
	}
	// publish invalid order
	err = sc.Publish(config.Stan.Channel, []byte(invalidOrder))
	if err != nil {
		t.Fatal(err)
	}
	// publish invalid order with too long Id
	err = sc.Publish(config.Stan.Channel, []byte(orderWithTooLongId))
	if err != nil {
		t.Fatal(err)
	}
	// publish valid order
	err = sc.Publish(config.Stan.Channel, []byte(validOrder))
	if err != nil {
		t.Fatal(err)
	}
	// publish duplicate valid order
	duplicateOrder := validOrder
	err = sc.Publish(config.Stan.Channel, []byte(duplicateOrder))
	if err != nil {
		t.Fatal(err)
	}
	// waiting for the messages to be processed
	time.Sleep(100 * time.Millisecond)

	router := api.SetupRouter()

	t.Run("Test store.ReadAllOrders func with non-empty db", func(t *testing.T) {
		orders, err := store.ReadAllOrders()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(len(cache.GetOrderIds()), len(orders))
		expected, _ := cache.GetOrder(validOrderId)
		actual := orders[validOrderId]
		assert.Equal(expected, actual)
	})

	t.Run("GET /api/order/", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/order/", nil)
		router.ServeHTTP(w, req)
		assert.Equal(http.StatusOK, w.Code)

		ids := new(api.OrderIds)
		err = json.Unmarshal(w.Body.Bytes(), ids)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(1, len(ids.Ids))
		assert.Equal(validOrderId, ids.Ids[0])
	})

	t.Run("GET /api/order/:id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/order/%s", validOrderId), nil)
		router.ServeHTTP(w, req)
		assert.Equal(http.StatusOK, w.Code)
		assert.JSONEq(validOrder, w.Body.String())

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/api/order/%s", invalidOrderId), nil)
		router.ServeHTTP(w, req)
		assert.Equal(http.StatusNotFound, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/api/order/%s", tooLongId), nil)
		router.ServeHTTP(w, req)
		assert.Equal(http.StatusNotFound, w.Code)
	})

	t.Run("GET /order/", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/order/", nil)
		router.ServeHTTP(w, req)
		assert.Equal(http.StatusOK, w.Code)
		assert.Contains(w.Body.String(), validOrderId)
	})

	t.Run("GET /order/:id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/order/%s", validOrderId), nil)
		router.ServeHTTP(w, req)
		assert.Equal(http.StatusOK, w.Code)
		assert.Contains(w.Body.String(), validOrderId)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/order/%s", invalidOrderId), nil)
		router.ServeHTTP(w, req)
		assert.Equal(http.StatusNotFound, w.Code)
	})
}

// func truncateAllTables(db *sql.DB) error {
// 	// https://stackoverflow.com/a/12082038/6475258
// 	truncateSql := `
// 		DO
// 		$do$
// 		BEGIN
// 			EXECUTE
// 			(SELECT 'TRUNCATE TABLE ' || string_agg(oid::regclass::text, ', ') || ' CASCADE'
// 			 FROM   pg_class
// 			 WHERE  relkind = 'r'  -- only tables
// 			 AND    relnamespace = 'public'::regnamespace
// 			);
// 		END
// 		$do$;
// 	`
// 	_, err := db.Exec(truncateSql)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func restoreDbFromTemplate(db *sql.DB) (*sql.DB, error) {
// 	err := db.Close()
// 	if err != nil {
// 		return nil, err
// 	}

// 	postgresDb := "postgres"
// 	postgresDbConnStr := fmt.Sprintf(
// 		"dbname=%s user=%s password=%s host=%s port=%s sslmode=%s",
// 		postgresDb,
// 		config.Postgres.User,
// 		config.Postgres.Password,
// 		config.Postgres.Host,
// 		config.Postgres.Port,
// 		config.Postgres.Sslmode,
// 	)
// 	pdb, err := store.MakeDB(postgresDbConnStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer pdb.Close()

// 	dropDbSql := fmt.Sprintf(`DROP DATABASE %s`, config.Postgres.Db)
// 	_, err = pdb.Exec(dropDbSql)
// 	if err != nil {
// 		return nil, err
// 	}

// 	createDbSql := fmt.Sprintf(`CREATE DATABASE %s TEMPLATE template_db`, config.Postgres.Db)
// 	_, err = pdb.Exec(createDbSql)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return store.MakeDB(config.Postgres.ConnStr)
// }
