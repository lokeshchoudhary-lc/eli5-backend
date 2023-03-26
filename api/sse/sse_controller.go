package sse

import (
	"bufio"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func LeaderboardSSE(c *fiber.Ctx) error {

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {

		clientConnection := make(chan []byte)
		log.Printf("client: %d", clientConnection)
		BroadcasterLB.newConnections <- clientConnection

		for {
			fmt.Fprintf(w, "data: %s\n\n", <-clientConnection)

			err := w.Flush()
			if err != nil {

				// Refreshing page in web browser will establish a new
				// SSE connection, but only (the last) one is alive, so
				// dead connections must be closed here.
				log.Printf("Error while flushing:' %v '. Closing http connection.\n", err)
				BroadcasterLB.closeConnection <- clientConnection
				break
			}
		}

	}))

	return nil
}

func BroadcastSSE(c *fiber.Ctx) error {
	BroadcasterLB.SendMessage()
	return c.SendStatus(200)
}
