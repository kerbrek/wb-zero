package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/kerbrek/wb-zero/app/cache"
	"github.com/kerbrek/wb-zero/app/model"
)

// GET /order/:id - renders order page
func getOrderPage(c *gin.Context) {
	var or = new(orderResource)
	if err := c.ShouldBindUri(or); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", "Bad Request")
		return
	}

	orderJson, ok := cache.GetOrder(or.Id)
	if !ok {
		c.HTML(http.StatusNotFound, "error.html", "Not Found")
		return
	}

	var order = new(model.Order)
	if err := json.Unmarshal(orderJson, order); err != nil {
		log.Errorf("getOrderPage: %v: %s", err, orderJson)
		c.HTML(http.StatusInternalServerError, "error.html", "Internal Server Error")
		return
	}

	c.HTML(http.StatusOK, "order.html", order)
}

// GET /order/ - renders page with list of order ids
func getOrderIdsPage(c *gin.Context) {
	ids := &OrderIds{Ids: cache.GetOrderIds()}
	c.HTML(http.StatusOK, "order_ids.html", ids)
}
