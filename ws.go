package main

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func getWSServer(channelName string) *socketio.Server {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	server.On("connection", func(so socketio.Socket) {
		log.Printf("New connection from %s ", so.Id())

		// listen from channel message from client and join client to the channel name
		so.On(channelName, func(channel string) {
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

	return server
}
