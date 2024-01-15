package main

import (
	"cache"
	"client/order"
	"db"
	"encoding/json"
	"log"
	"time"
	"web"

	"github.com/nats-io/stan.go"
)

var (
	clusterID string = "test-cluster"
	clientID  string = "1"
	channel   string = "Work"
	memory    *cache.Cache
)

func cacheRecovery(memory *cache.Cache) {
	orders := db.GetAllOrders()
	if len(orders) > 0 {
		for i := 0; i < len(orders); i++ {
			memory.Set(orders[i].OrderUID, &orders[i])
		}
	}
}

func main() {
	memory = cache.NewCache()
	go cacheRecovery(memory)
	go web.HandleReq(memory)
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatalf("error connection : %v", err)
	}
	defer sc.Close()

	messageHandler := func(msg *stan.Msg) {
		var dt = order.Order{}
		err := json.Unmarshal(msg.Data, &dt)
		if err != nil {
			log.Fatal(err)
		}
		db.SaveInDB(&dt)
		record := db.GetRecord(dt.OrderUID)
		cache.SaveInMemory(memory, &record)

	}

	subscription, err := sc.Subscribe(channel, messageHandler, stan.DurableName("durable-name"))
	if err != nil {
		log.Fatal(err)
	}
	defer subscription.Close()
	for {
		time.Sleep(time.Second)
	}
}
