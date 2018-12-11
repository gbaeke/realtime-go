package main

import (
	"log"
	"net/http"

	"github.com/go-redis/redis"
	socketio "github.com/googollee/go-socket.io"
)

func main() {
	// redis connection
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// subscribe to all channels
	pubsub := client.PSubscribe("*")

	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	// messages received on a Go channel
	ch := pubsub.Channel()

	// ping Redis server
	pong, err := client.Ping().Result()
	log.Println(pong, err)

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Printf("New connection from %s ", so.Id())

		// listen from channel message from client and join client to the channel name
		so.On("channel", func(channel string) {
			log.Printf("%s joins channel %s\n", so.Id(), channel)
			so.Join(channel)
			so.BroadcastTo(channel, "hello")
		})

		so.On("disconnection", func() {
			log.Printf("disconnect from %s\n", so.Id())
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	// Consume messages from Redis
	go func(srv *socketio.Server) {
		for msg := range ch {
			log.Println(msg.Channel, msg.Payload)
			srv.BroadcastTo(msg.Channel, "message", msg.Payload)
		}
	}(server)

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./assets")))
	log.Println("Serving on localhost:8888...")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
