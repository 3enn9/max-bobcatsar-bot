package max

import (
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"net/http"
)

func WebhookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var maxBotApi schemes.Update
	}
}
