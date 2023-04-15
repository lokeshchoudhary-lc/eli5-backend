package sse

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
)

type lbBroadcaster struct {
	broadcaster     chan []byte
	newConnections  chan chan []byte
	closeConnection chan chan []byte
	connections     map[chan []byte]bool
}

func newLBBroadcaster() *lbBroadcaster {
	return &lbBroadcaster{
		broadcaster:     make(chan []byte),
		newConnections:  make(chan chan []byte),
		closeConnection: make(chan chan []byte),
		connections:     make(map[chan []byte]bool),
	}
}

func newLBServer() *lbBroadcaster {
	broadcaster := newLBBroadcaster()
	go broadcaster.lbListen()
	return broadcaster
}

func (b *lbBroadcaster) lbListen() {
	for {
		select {
		case connection := <-b.newConnections:
			b.connections[connection] = true
			log.Printf("New Connection. %d registered Connection", len(b.connections))

		case connection := <-b.closeConnection:
			delete(b.connections, connection)
			log.Printf("Removed Connection. %d registered Connection", len(b.connections))

		case message := <-b.broadcaster:
			log.Printf("message: %s", message)
			for connection := range b.connections {
				log.Printf("client send: %d", connection)
				connection <- message
			}
		}
	}
}

type user struct {
	UserId       string `json:"userId"`
	TotalAnswers int64  `json:"totalAnswers,omitempty"`
	TotalLikes   int64  `json:"totalLikes,omitempty"`
	Rank         int64  `json:"rank,omitempty"`
}

func (b *lbBroadcaster) SendMessage() {

	user := &user{UserId: strconv.Itoa((rand.Intn(100))), TotalAnswers: 10, TotalLikes: 5, Rank: 12}

	data, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	b.broadcaster <- []byte(string(data))
	// b.broadcaster <- byte
	// b.broadcaster <- []byte(strconv.Itoa((rand.Intn(10))))
}

// func startPinging(b *lbBroadcaster) {
// 	ticker := time.NewTicker(2 * time.Second)
// 	quit := make(chan struct{})

// 	go func() {
// 		for {
// 			select {
// 			case <-ticker.C:
// 				user := &user{UserId: strconv.Itoa((rand.Intn(100))), TotalAnswers: 10, TotalLikes: 5, Rank: 12}

// 				data, err := json.Marshal(user)
// 				if err != nil {
// 					panic(err)
// 				}
// 				b.broadcaster <- []byte(string(data))
// 				// b.broadcaster <- []byte(strconv.Itoa((rand.Intn(10))))
// 			case <-quit:
// 				ticker.Stop()
// 				return
// 			}
// 		}
// 	}()
// }

var BroadcasterLB *lbBroadcaster

func init() {

	BroadcasterLB = newLBServer()
	// startPinging(BroadcasterLB)
}
