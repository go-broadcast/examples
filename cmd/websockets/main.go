package main

import (
	"log"

	"github.com/go-broadcast/broadcast"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("./static/websockets/index.html")
	})

	app.Static("/static/websockets", "../../static/websockets")

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	bcast, err := broadcast.New()
	if err != nil {
		log.Fatal(err)
	}

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		user := c.Params("id")

		sub := bcast.Subscribe(func(data interface{}) {
			c.WriteJSON(data)
		})
		bcast.JoinRoom(sub, "user-"+user, "chat-room")

		var (
			msg []byte
			err error
		)
		for {
			if _, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}

			bcast.ToRoom(string(msg), "chat-room", "user-"+user)
		}

		bcast.Unsubscribe(sub)
	}))

	app.Listen(":5200")
}
