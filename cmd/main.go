package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()

	host := os.Getenv("HOST")
	token := os.Getenv("TOKEN")
	secret := os.Getenv("SECRET")

	api, _ := maxbot.New(token)

	errChan := api.GetErrors()
	go func() {
		for errMessage := range errChan {
			log.Println(errMessage) // use your favorite logger
		}
	}()

	// Some methods demo:
	info, err := api.Bots.GetBot(ctx)
	fmt.Printf("Get me: %#v %#v", info, err)

	subs, _ := api.Subscriptions.GetSubscriptions(ctx)
	for _, s := range subs.Subscriptions {
		_, _ = api.Subscriptions.Unsubscribe(ctx, s.Url)
	}

	subscriptionResp, err := api.Subscriptions.Subscribe(ctx, host+"/webhook", []string{}, secret)
	log.Printf("Subscription: %#v %#v", subscriptionResp, err)

	ch := make(chan schemes.UpdateInterface)

	http.HandleFunc("/webhook", api.GetHandler(ch))
	go func() {
		for {
			update := <-ch
			log.Printf("Received: %#v", update)
			switch upd := update.(type) {
			case *schemes.MessageCreatedUpdate:
				message := maxbot.NewMessage().
					SetUser(upd.Message.Sender.UserId).
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
