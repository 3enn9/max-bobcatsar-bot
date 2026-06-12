package main

import (
	"bobcatsar-max-bot/internal/config"
	"bobcatsar-max-bot/internal/db"
	"bobcatsar-max-bot/internal/grpc/accountant"
	max2 "bobcatsar-max-bot/internal/max"
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

	pool, err := db.ConnectionDB(cfg)
	if err != nil {
		log.Fatalf("error create db pool %v", err)
	}
	rep := db.NewRepository(pool)

	api, err := maxbot.New(cfg.Token)
	if err != nil {
		log.Fatalf("failed to create api: %v", err)
	}

	accountantClient, err := accountant.NewClient("python-service:50051")

	if err != nil {
		log.Fatal(err)
	}
	defer accountantClient.Close()

	updateCh := make(chan schemes.UpdateInterface, 3)
	maxService := max2.NewMaxService(rep, api, updateCh, accountantClient)

	errChan := api.GetErrors()
	go func() {
		for errMessage := range errChan {
			log.Println(errMessage)
		}
	}()

	info, err := api.Bots.GetBot(ctx)
	fmt.Printf("Get me: %#v %#v", info, err)

	http.HandleFunc("/webhook", api.GetHandler(updateCh))

	go maxService.ListenUpdates(ctx)

	_ = http.ListenAndServe(":8080", nil)
}
