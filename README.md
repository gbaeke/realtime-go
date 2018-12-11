# Real-time socket.io app in Go

Real-time app that uses go-socket.io (github.com/googollee/go-socket.io)

This is a simple example without authentication or other more advanced features. It could be used for demo purposes in conjuction with an IoT device that sends data you want to display in the browser in real-time.

It connects to a local Redis host on port 6379 and subscribes to all Redis channels.

Whenever a message is received on a Redis channel, it broadcasts the message payload to all connected sockets that have joined a *room* that equals the channel name.

The assets folder contains a simple web app written in Vue that uses the socket.io client to join the room device01. Whenever you publish a messages to a Redis channel with that name, the payload will be displayed.

You can also run the solution is Azure ACI with the following command:

az group create -g realtime-rg -l westeurope

az container create --resource-group realtime-rg --file multi.yaml

When the container group is created, go the the IP of the container group (e.g. http://ip:8888)
