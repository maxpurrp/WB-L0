package main

import (
	"client/order"
	"encoding/json"
	"log"

	"github.com/nats-io/stan.go"
)

var (
	clusterID string = "test-cluster"
	clientID  string = "2"
	channel   string = "Work"
)

func main() {
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatal(err)
	}

	defer sc.Close()
	order, err := order.GetOrder()
	if err != nil {
		log.Fatal(err)
	}

	jsonOrder, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}

	// send json to channes
	err = sc.Publish(channel, jsonOrder)
	if err != nil {
		log.Fatal(err)
	}
}
