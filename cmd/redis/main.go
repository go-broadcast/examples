package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-broadcast/broadcast"
	"github.com/go-broadcast/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	port := flag.Int("port", 0, "service port")
	flag.Parse()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("./static/redis/index.html")
	})

	app.Static("/static/redis", "../../static/redis")

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	dispatcher, err := redis.New()
	if err != nil {
		log.Fatal(err)
	}

	broadcaster, close, err := broadcast.New(
		broadcast.WithDispatcher(dispatcher),
	)
	if err != nil {
		log.Fatal(err)
	}

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		user := c.Params("id")

		subscription := broadcaster.Subscribe(func(data interface{}) {
			c.WriteJSON(data)
		})
		broadcaster.JoinRoom(subscription, "user-"+user, "chat-room")

		var (
			msg []byte
			err error
		)
		for {
			if _, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}

			broadcaster.ToRoom(string(msg), "chat-room", "user-"+user)
		}

		broadcaster.Unsubscribe(subscription)
	}))

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		app.Shutdown()
	}()

	app.Listen(":" + strconv.Itoa(*port))
	log.Println("stopped fiber app")

	close()
	<-broadcaster.Done()
	log.Println("stopped broadcaster")
}
