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
		"/salary": m.PrePaymentMessage,
	}
	return m
}

func (ms *MaxService) PrePaymentMessage(upd *schemes.MessageCreatedUpdate) {
	salaryStr := strings.Split(upd.Message.Body.Text, " ")[1]
	salaryFloat, err := strconv.ParseFloat(salaryStr, 64)
	if err != nil {
		log.Println("error convert string to float")
		return
	}
	db.AddPrePayment(ms.pool, " ", salaryFloat, upd.Message.Recipient.ChatId)
}
