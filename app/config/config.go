package config

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type postgresConf struct {
	Db       string
	User     string
	Password string
	Host     string
	Port     string
	Sslmode  string
	ConnStr  string
}

type stanConf struct {
	Host        string
	Port        string
	ConnStr     string
	ClusterID   string
	ClientID    string
	Channel     string
	DurableName string
}

type appConf struct {
	GinMode string
	Host    string
	Port    string
	Addr    string
}

func getenvOrFatal(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("getenvOrFatal: environment variable '%s' is empty", key)
	}
	return value
}

func getenvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

var Postgres postgresConf
var Stan stanConf
var App appConf

func init() {
	Postgres.Db = getenvOrFatal("POSTGRES_DB")
	Postgres.User = getenvOrFatal("POSTGRES_USER")
	Postgres.Password = getenvOrFatal("POSTGRES_PASSWORD")
	Postgres.Host = getenvOrFatal("POSTGRES_HOST")
	Postgres.Port = getenvOrFatal("POSTGRES_PORT")
	Postgres.Sslmode = getenvOrFatal("POSTGRES_SSLMODE")
	Postgres.ConnStr = fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%s sslmode=%s",
		Postgres.Db,
		Postgres.User,
		Postgres.Password,
		Postgres.Host,
		Postgres.Port,
		Postgres.Sslmode,
	)

	Stan.Host = getenvOrFatal("STAN_HOST")
	Stan.Port = getenvOrFatal("STAN_PORT")
	Stan.ConnStr = fmt.Sprintf(
		"nats://%s:%s",
		Stan.Host,
		Stan.Port,
	)
	Stan.ClusterID = getenvOrFatal("STAN_CLUSTER_ID")
	Stan.ClientID = getenvOrFatal("STAN_CLIENT_ID")
	Stan.Channel = getenvOrFatal("STAN_CHANNEL")
	Stan.DurableName = getenvOrFatal("STAN_DURABLE_NAME")

	App.GinMode = getenvOrDefault(gin.EnvGinMode, gin.ReleaseMode)
	App.Host = getenvOrFatal("APP_HOST")
	App.Port = getenvOrFatal("APP_PORT")
	App.Addr = fmt.Sprintf(
		"%s:%s",
		App.Host,
		App.Port,
	)
}
