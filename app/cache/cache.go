package cache

import (
	"sync"
)

var cache = struct {
	sync.RWMutex
	m map[string][]byte
}{}

func Init(orderJsons map[string][]byte) {
	cache.Lock()
	cache.m = orderJsons
	cache.Unlock()
}

func SetOrder(orderId string, orderJson []byte) {
	cache.Lock()
	cache.m[orderId] = orderJson
	cache.Unlock()
}

func GetOrder(orderId string) ([]byte, bool) {
	cache.RLock()
	orderJson, ok := cache.m[orderId]
	cache.RUnlock()
	return orderJson, ok
}

func GetOrderIds() []string {
	cache.RLock()
	orderIds := make([]string, 0, len(cache.m))
	for id := range cache.m {
		orderIds = append(orderIds, id)
	}
	cache.RUnlock()
	return orderIds
}

// // Cache implemented using sync.Map
// var cache = sync.Map{}

// func Init(orderJsons map[string][]byte) {
// 	for id, oj := range orderJsons {
// 		cache.Store(id, oj)
// 	}
// }

// func SetOrder(orderId string, orderJson []byte) {
// 	cache.Store(orderId, orderJson)
// }

// func GetOrder(orderId string) ([]byte, bool) {
// 	orderJson, ok := cache.Load(orderId)
// 	if !ok {
// 		return nil, ok
// 	}
// 	return orderJson.([]byte), ok
// }

// func GetOrderIds() []string {
// 	ids := make([]string, 0)
// 	cache.Range(func(id, _ any) bool {
// 		ids = append(ids, id.(string))
// 		return true
// 	})
// 	return ids
// }
