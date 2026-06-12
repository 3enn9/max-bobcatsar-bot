package max

import (
	"bobcatsar-max-bot/internal/db"
	"bobcatsar-max-bot/internal/grpc/accountant"
	"bobcatsar-max-bot/internal/grpc/gen"
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
	Client   *accountant.Client
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

func NewMaxService(rep *db.Repository, maxApi *maxbot.Api, updateCh <-chan schemes.UpdateInterface, client *accountant.Client) *MaxService {
	m := &MaxService{rep: rep, maxApi: maxApi, updateCh: updateCh, Client: client}
	m.Commands = map[string]CommandHandler{
		"/salary": m.PrePaymentCommand,
		"/show":   m.ShowPrePayments,
		"/add":    m.handleAdd,
		"/send":   m.SendMessageAccountant,
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

func (ms *MaxService) handleAdd(upd *schemes.MessageCreatedUpdate) *maxbot.Message {

	text := upd.Message.Body.Text
	operationArray := strings.Split(text, "\n")
	chatId := upd.Message.Recipient.ChatId

	var operations []string
	var errorsList []string
	var totalAmount float64

	for i, operation := range operationArray {
		parts := strings.Fields(operation)

		if len(parts) < 2 {
			errorsList = append(errorsList, operation+" (неверный формат)")
			continue
		}

		if i == 0 {
			parts = parts[1:]
		}

		amountStr := parts[len(parts)-1]
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			errorsList = append(errorsList, operation+" (неверная сумма)")
			continue
		}
		description := strings.Join(parts[:len(parts)-1], " ")
		operationType := "withdraw"
		if amount >= 0 {
			operationType = "deposit"
		}

		payment := db.Payment{Description: description, Operation: operationType, TelegramGroupID: chatId, Amount: amount}

		err = ms.rep.AddPayment(payment)
		if err != nil {
			log.Printf("Не удалось добавить операцию %v", err)
			continue
		}

		operations = append(operations, fmt.Sprintf("• %s: %.2f", description, amount))
		totalAmount += amount

	}
	balance, err := ms.rep.GetBalance(chatId)
	if err != nil {
		log.Printf("Ошибка получения баланса: %v", err)
	}

	answer := fmt.Sprintf(
		"📊 Операции:\n%s\n\n💰 Итого: %.2f\n🏦 Касса: %.2f",
		strings.Join(operations, "\n"),
		totalAmount,
		balance,
	)

	if len(errorsList) > 0 {
		answer += "\n\n⚠ Пропущены:\n" + strings.Join(errorsList, "\n")
	}
	msg := maxbot.NewMessage().SetChat(chatId).SetText(answer)

	return msg
}

func (ms *MaxService) SendMessageAccountant(upd *schemes.MessageCreatedUpdate) *maxbot.Message {
	text := strings.Fields(upd.Message.Body.Text)
	chatId := upd.Message.Recipient.ChatId
	var msg *maxbot.Message
	if len(text) < 2 {
		log.Println("Не верный формат ввода команды")
		return nil
	}

	request := &gen.SendAccountantRequest{Message: "test"} // заглушка

	response, err := ms.Client.SendAccountant(context.Background(), request)
	if err != nil {
		fmt.Printf("error %v", err)
	}

	if response.Success {
		msg = maxbot.NewMessage().SetChat(chatId).SetText("Отправлено")
	} else {
		msg = maxbot.NewMessage().SetChat(chatId).SetText("Не получилось отправить")
	}
	return msg
}
