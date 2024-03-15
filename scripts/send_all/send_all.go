package main

import (
	"context"
	"fmt"
	"quree/config"
	"quree/internal/models"
	"quree/internal/nats"
	"quree/internal/pg"
	"sync"
)

var ctx = context.Background()
var db = pg.DB
var ncAdmin *nats.NatsClient = nats.Init(nats.NatsSettings{
	Ctx: ctx,
	URL: config.NATS_URL,
})

var ncUser *nats.NatsClient = nats.Init(nats.NatsSettings{
	Ctx: ctx,
	URL: config.NATS_URL,
})

func main() {

	wg := &sync.WaitGroup{}
	ncAdmin.UsePublishSubject(config.NATS_ADMIN_MESSAGES_SUBJECT)
	ncUser.UsePublishSubject(config.NATS_USER_MESSAGES_SUBJECT)

	adminIds := db.GetAdminChatIDs()
	fmt.Println("Admin: ", len(adminIds))

	// targets := []string{"306562182", "6696815781"}

	text := "Всем добрый день! Для проверки допуска на финальное событие была добавлена кнопка \"Проверить проходку\".\n\n"
	text += "Если у вас кнопка не появилась, введите команду /start"

	message := &models.SendableMessage{
		Text: &text,
	}

	wg.Add(1)
	sendAll(&adminIds, message, ncAdmin)
	wg.Wait()

}

// write a function that sends a message to all admins

func sendAll(targetChatIDs *[]string, message *models.SendableMessage, nc *nats.NatsClient) {
	for _, chatID := range *targetChatIDs {
		message.Recipient = &models.Recipient{
			ChatID: chatID,
		}
		nc.Publish(message)
	}
	fmt.Println("All sent!")
}
