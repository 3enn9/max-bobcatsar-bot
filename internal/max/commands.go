package max

import (
	"bobcatsar-max-bot/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"log"
	"strconv"
	"strings"
)

type CommandHandler func(update *schemes.MessageCreatedUpdate) *maxbot.Message

type MaxService struct {
	pool     *pgxpool.Pool
	Commands map[string]CommandHandler
}

func NewMaxService(pool *pgxpool.Pool) *MaxService {
	m := &MaxService{pool: pool}
	m.Commands = map[string]CommandHandler{
		"/salary": m.PrePaymentCommand,
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
	err = db.AddPrePayment(ms.pool, " ", salary, upd.Message.Recipient.ChatId)

	if err != nil {
		log.Printf("Не удалось добавить запись в бд ошибка: %v\n", err)
		return nil
	}
	msg := maxbot.NewMessage().
		SetChat(upd.Message.Recipient.ChatId).
		SetText("Аванс успешно добавлен")

	return msg

}
