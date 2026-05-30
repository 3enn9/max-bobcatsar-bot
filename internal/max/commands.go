package max

import (
	"bobcatsar-max-bot/internal/db"
	"context"
	"fmt"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"log"
	"strconv"
	"strings"
)

type CommandHandler func(update *schemes.MessageCreatedUpdate) *maxbot.Message

type MaxService struct {
	rep      *db.Repository
	maxApi   *maxbot.Api
	Commands map[string]CommandHandler
	updateCh <-chan schemes.UpdateInterface
}

func (ms *MaxService) ListenUpdates(ctx context.Context) {
	for {
		select {
		case update, ok := <-ms.updateCh:
			if !ok {
				return
			}
			log.Printf("Received: %#v", update)
			switch upd := update.(type) {
			case *schemes.MessageCreatedUpdate:
				command := upd.Message.Body.Text
				command = strings.Split(command, " ")[0]

				if someFunc, ok := ms.Commands[command]; ok {
					go func(upd *schemes.MessageCreatedUpdate) {
						msg := someFunc(upd)
						err := ms.maxApi.Messages.Send(ctx, msg)
						if err != nil {
							fmt.Printf("error send message %v", err)
						}
					}(upd)
				}
			default:
				log.Printf("Unknown type: %#v", upd)
			}
		case <-ctx.Done():
			return
		}
	}
}

func NewMaxService(rep *db.Repository, maxApi *maxbot.Api, updateCh <-chan schemes.UpdateInterface) *MaxService {
	m := &MaxService{rep: rep, maxApi: maxApi, updateCh: updateCh}
	m.Commands = map[string]CommandHandler{
		"/salary": m.PrePaymentCommand,
		"/show":   m.ShowPrePayments,
	}
	return m
}

func (ms *MaxService) PrePaymentCommand(upd *schemes.MessageCreatedUpdate) *maxbot.Message {
	text := strings.Fields(upd.Message.Body.Text)
	if len(text) != 2 {
		log.Println("Не верный формат ввода команды")
		return nil
	}
	salary, err := strconv.ParseFloat(text[1], 64)
	if err != nil {
		log.Println("error convert string to float")
		return nil
	}
	err = ms.rep.AddPrePayment(" ", salary, upd.Message.Recipient.ChatId)

	if err != nil {
		log.Printf("Не удалось добавить запись в бд ошибка: %v\n", err)
		return nil
	}
	msg := maxbot.NewMessage().
		SetChat(upd.Message.Recipient.ChatId).
		SetText("Аванс успешно добавлен")

	return msg

}

func (ms *MaxService) ShowPrePayments(upd *schemes.MessageCreatedUpdate) *maxbot.Message {
	chatID := upd.Message.Recipient.ChatId
	err, text := ms.rep.PrePayments(chatID)
	if err != nil {
		log.Printf("Не удалось показать авансы")
		return nil
	}
	msg := maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetText(text)

	return msg
}
