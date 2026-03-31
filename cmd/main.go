package main

import (
	"bobcatsar-max-bot/internal/config"
	"bobcatsar-max-bot/internal/db"
	"context"
	"fmt"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()
	cfg := config.NewConfig()

	db.ConnectionDB(cfg)

	api, _ := maxbot.New(cfg.Token)

	errChan := api.GetErrors()
	go func() {
		for errMessage := range errChan {
			log.Println(errMessage) // use your favorite logger
		}
	}()

	// Some methods demo:
	info, err := api.Bots.GetBot(ctx)
	fmt.Printf("Get me: %#v %#v", info, err)

	ch := make(chan schemes.UpdateInterface)

	http.HandleFunc("/webhook", api.GetHandler(ch))
	go func() {
		for {
			update := <-ch
			log.Printf("Received: %#v", update)
			switch upd := update.(type) {
			case *schemes.MessageCreatedUpdate:
				message := maxbot.NewMessage().
					SetUser(upd.Message.Recipient.ChatId).
					SetText(fmt.Sprintf("Hello, %s! Your message: %s", upd.Message.Sender.Name, upd.Message.Body.Text))

				err = api.Messages.Send(ctx, message)
				log.Printf("Answer: %#v", err)
			default:
				log.Printf("Unknown type: %#v", upd)
			}
		}
	}()

	_ = http.ListenAndServe(":8080", nil)
}
