package sse

// import (
// 	"log"
// 	"math/rand"
// 	"strconv"
// 	"time"
// )

// type userBroadcaster struct {
// 	broadcaster     chan []byte
// 	newConnections  chan chan []byte
// 	closeConnection chan chan []byte
// 	connections     map[chan []byte]bool
// }

// func newUserBroadcaster() *userBroadcaster {
// 	return &userBroadcaster{
// 		broadcaster:     make(chan []byte),
// 		newConnections:  make(chan chan []byte),
// 		closeConnection: make(chan chan []byte),
// 		connections:     make(map[chan []byte]bool),
// 	}
// }

// type userServer struct {
// 	broadcaster *userBroadcaster
// }

// func newUserServer() *userServer {
// 	broadcaster := newUserBroadcaster()
// 	server := &userServer{
// 		broadcaster,
// 	}
// 	go broadcaster.userListen()
// 	return server
// }

// func startPinging(b *lbBroadcaster) {
// 	ticker := time.NewTicker(2 * time.Second)
// 	quit := make(chan struct{})

// 	go func() {
// 		for {
// 			select {
// 			case <-ticker.C:
// 				b.broadcaster <- []byte(strconv.Itoa((rand.Intn(10))))
// 			case <-quit:
// 				ticker.Stop()
// 				return
// 			}
// 		}
// 	}()
// }

// func (b *userBroadcaster) userListen() {
// 	for {
// 		select {
// 		case connection := <-b.newConnections:
// 			b.connections[connection] = true
// 			log.Printf("New Connection. %d registered Connection", len(b.connections))

// 		case connection := <-b.closeConnection:
// 			delete(b.connections, connection)
// 			log.Printf("Removed Connection. %d registered Connection", len(b.connections))

// 		case message := <-b.broadcaster:
// 			log.Printf("message: %s", message)
// 			for connection := range b.connections {
// 				log.Printf("client send: %d", connection)
// 				connection <- message
// 			}
// 		}
// 	}
// }
