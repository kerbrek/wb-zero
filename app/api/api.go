package api

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kerbrek/wb-zero/app/cache"
	"github.com/kerbrek/wb-zero/app/config"
)

func formatAsCurrency(currency string) string {
	switch currency {
	case "USD":
		return "$"
	default:
		return currency
	}
}

func formatAsPrice(price uint64) string {
	basicMonetaryUnits := price / 100   // dollar, euro, ruble, etc.
	oneHundredthSubunits := price % 100 // cent, kopek, etc.
	return fmt.Sprintf("%d.%02d", basicMonetaryUnits, oneHundredthSubunits)
}

func formatAsDate(unixTime int64) string {
	year, month, day := time.Unix(unixTime, 0).Date()
	return fmt.Sprintf("%d %s %d", day, month, year)
}

func SetupRouter() *gin.Engine {
	gin.SetMode(config.App.GinMode)
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"formatAsDate":     formatAsDate,
		"formatAsPrice":    formatAsPrice,
		"formatAsCurrency": formatAsCurrency,
	})
	r.LoadHTMLGlob("templates/*")

	r.GET("/api/order/", getOrderIds)
	r.GET("/api/order/:id", getOrder)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/order/")
	})
	r.GET("/order/", getOrderIdsPage)
	r.GET("/order/:id", getOrderPage)
	return r
}

type orderResource struct {
	Id string `uri:"id" binding:"required"`
}

// GET /api/order/:id - returns order
func getOrder(c *gin.Context) {
	var or = new(orderResource)
	if err := c.ShouldBindUri(or); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}

	orderJson, ok := cache.GetOrder(or.Id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Order does not exist."})
		return
	}

	c.Data(http.StatusOK, "application/json", orderJson)
}

type OrderIds struct {
	Ids []string `json:"order_ids"`
}

// GET /api/order/ - returns list of order ids
func getOrderIds(c *gin.Context) {
	ids := &OrderIds{Ids: cache.GetOrderIds()}
	c.JSON(http.StatusOK, ids)
}
