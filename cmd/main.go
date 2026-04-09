package main

import (
	"bobcatsar-max-bot/internal/config"
	"bobcatsar-max-bot/internal/db"
	"bobcatsar-max-bot/internal/max"
	"context"
	"fmt"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()
	cfg := config.NewConfig()

	pool, err := db.ConnectionDB(cfg)

	if err != nil {
		log.Fatalf("error create db pool %v", err)
	}

	maxService := max.NewMaxService(pool)
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
				command := upd.Message.Body.Text
				command = strings.Split(command, " ")[0]

				if someFunc, ok := maxService.Commands[command]; ok {
					msg := someFunc(upd)
					err = api.Messages.Send(context.Background(), msg)
					if err != nil {
						log.Printf("Ошибка отправки сообщения %v", err)
					}
				}
			default:
				log.Printf("Unknown type: %#v", upd)
			}
		}
	}()

	_ = http.ListenAndServe(":8080", nil)
}
