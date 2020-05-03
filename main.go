package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	socketio "github.com/googollee/go-socket.io"
	consulapi "github.com/hashicorp/consul/api"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {
	// consul
	config := consulapi.DefaultConfig()
	config.Address = "http://consul.service.consul:8500"
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	kv := consul.KV()
	redisPattern, _, err := kv.Get("REDISPATTERN", nil)
	if err != nil || redisPattern == nil {
		log.Fatalf("Could not get REDISPATTERN: %s", err)
	}
	fmt.Printf("KV: %v %s\n", redisPattern.Key, redisPattern.Value)

	// redis connection
	client := redis.NewClient(&redis.Options{
		Addr: getEnv("REDISHOST", "localhost:6379"),
	})

	// subscribe to all channels
	pubsub := client.PSubscribe(string(redisPattern.Value))

	_, err = pubsub.Receive()
	if err != nil {
		panic(err)
	}

	// messages received on a Go channel
	ch := pubsub.Channel()

	// ping Redis server
	pong, err := client.Ping().Result()
	log.Println(pong, err)

	server := getWSServer("channel")
	if server == nil {
		log.Fatalln("Could not create WebSockets server")
	}

	// Consume messages from Redis
	go func(srv *socketio.Server) {
		for msg := range ch {
			log.Println(msg.Channel, msg.Payload)
			srv.BroadcastTo(msg.Channel, "message", msg.Payload)
		}
	}(server)

	mux := http.NewServeMux()
	mux.Handle("/socket.io/", server)
	//mux.Handle("/", http.FileServer(http.Dir("./assets")))
	mux.HandleFunc("/", handle)

	http.ListenAndServe(":8080", mux)

}

func handle(w http.ResponseWriter, req *http.Request) {

	http.ServeFile(w, req, "./assets")

}
