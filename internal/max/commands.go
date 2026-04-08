package max

import (
	"bobcatsar-max-bot/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"log"
	"strconv"
	"strings"
)

type CommandHandler func(update *schemes.MessageCreatedUpdate)

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

func (ms *MaxService) PrePaymentCommand(upd *schemes.MessageCreatedUpdate) {
	text := strings.Split(upd.Message.Body.Text, " ")
	if len(text) <= 1 {
		log.Println("Не верный формат ввода команды")
		return
	}
	salary, err := strconv.ParseFloat(text[1], 64)
	if err != nil {
		log.Println("error convert string to float")
		return
	}
	db.AddPrePayment(ms.pool, " ", salary, upd.Message.Recipient.ChatId)
}
