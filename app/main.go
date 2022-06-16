package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/kerbrek/wb-zero/app/api"
	"github.com/kerbrek/wb-zero/app/broker"
	"github.com/kerbrek/wb-zero/app/cache"
	"github.com/kerbrek/wb-zero/app/config"
	"github.com/kerbrek/wb-zero/app/store"
)

func main() {
	db, err := store.MakeDB(config.Postgres.ConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	orderJsons, err := store.ReadAllOrders()
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	cache.Init(orderJsons)

	_, sc, subErr, connErr := broker.MakeStan()
	if connErr != nil {
		db.Close()
		log.Fatal(connErr)
	}
	if subErr != nil {
		sc.Close()
		db.Close()
		log.Fatal(subErr)
	}
	defer sc.Close()

	router := api.SetupRouter()
	err = router.Run(config.App.Addr)
	if err != nil {
		sc.Close()
		db.Close()
		log.Fatal(err)
	}
}
