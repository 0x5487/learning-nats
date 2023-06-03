package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {

	nc, err := nats.Connect("nats://host.docker.internal:4222")
	if err != nil {
		panic(err)
	}
	defer nc.Drain()

	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	cfg := nats.StreamConfig{
		Name:     "EVENTS",
		Storage:  nats.MemoryStorage,
		Subjects: []string{"events.>"},
	}

	_, err = js.AddStream(&cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println("created the stream")

	start := time.Now()

	for i := 0; i < 10000; i++ {
		_, err := js.PublishAsync("events.created_order", nil)
		if err != nil {
			break
		}
	}

	select {
	case <-js.PublishAsyncComplete():
		fmt.Println("published 6 messages")
	case <-time.After(time.Second):
		log.Fatal("publish took too long")
	}

	elapsed := time.Since(start)
	log.Printf("end took %s", elapsed)

	printStreamState(js, cfg.Name)
}

func printStreamState(js nats.JetStreamContext, name string) {
	info, err := js.StreamInfo(name)
	if err != nil {
		panic(err)
	}
	b, _ := json.MarshalIndent(info.State, "", " ")
	fmt.Println("inspecting stream info")
	fmt.Println(string(b))
}
