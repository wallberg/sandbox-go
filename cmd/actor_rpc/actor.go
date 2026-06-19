package main

import (
	"fmt"
	"net/rpc"
)

// The RPC Actor example from
// https://medium.com/@joao_vaz/actor-model-for-concurrent-systems-an-introduction-in-go-75fd25f2f04e

type Message struct {
	To   string
	From string
	Body string
}

type Actor struct {
	Name     string
	Messages []Message
}

func (a *Actor) SendMessage(m Message) {
	a.Messages = append(a.Messages, m)
}

func (a *Actor) ProcessMessages() {
	for _, m := range a.Messages {
		fmt.Printf("Actor %s received message: %s\n", a.Name, m.Body)
	}
	a.Messages = nil
}

type ActorManager struct {
	Actors map[string]*Actor
}

func NewActorManager() *ActorManager {
	return &ActorManager{Actors: make(map[string]*Actor)}
}

func (s *ActorManager) RegisterActor(name string) {
	s.Actors[name] = &Actor{Name: name}
}

func (s *ActorManager) SendMessage(m Message) {
	if actor, ok := s.Actors[m.To]; ok {
		actor.SendMessage(m)
	} else {
		// send message to remote actor
		client, err := rpc.Dial("tcp", m.To)
		if err != nil {
			fmt.Printf("Error connecting to remote actor %s: %s\n", m.To, err)
			return
		}
		defer client.Close()
		var reply string
		err = client.Call("RemoteActor.ReceiveMessage", m, &reply)
		if err != nil {
			fmt.Printf("Error sending message to remote actor %s: %s\n", m.To, err)
			return
		}
	}
}

type RemoteActor struct{}

func (a *RemoteActor) ReceiveMessage(m Message, reply *string) error {
	fmt.Printf("Remote actor %s received message: %s\n", m.To, m.Body)
	*reply = "Message received"
	return nil
}
