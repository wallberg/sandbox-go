package main

import (
	"fmt"
	"net"
	"net/rpc"
)

// The RPC Actor example from
// https://medium.com/@joao_vaz/actor-model-for-concurrent-systems-an-introduction-in-go-75fd25f2f04e

func main() {
	// start the remote actor (simulate a separate process)
	actor := new(RemoteActor)
	rpc.Register(actor)
	l, e := net.Listen("tcp", "localhost:1234")
	if e != nil {
		fmt.Printf("Error listening: %s\n", e)
		return
	}
	defer l.Close()

	fmt.Println("Remote actor listening on localhost:1234")

	go rpc.Accept(l)

	// create local actor system
	actorManager := NewActorManager()

	// register local actors
	actorManager.RegisterActor("actor1")
	actorManager.RegisterActor("actor2")

	// send message to local actor
	actorManager.SendMessage(Message{
		To:   "actor1",
		From: "actor2",
		Body: "Hello from actor2 to actor1",
	})

	// send message to remote actor
	actorManager.SendMessage(Message{
		To:   "localhost:1234",
		From: "actor1",
		Body: "Hello from actor1 to remote actor on localhost:1234",
	})

	// process messages in all actors
	for _, actor := range actorManager.Actors {
		actor.ProcessMessages()
	}
}
