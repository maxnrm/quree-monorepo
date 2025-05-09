package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"quree/config"
	"quree/internal/models"
	"quree/internal/nats"
	"sync"
	"time"

	"gopkg.in/telebot.v3"
)

type User struct {
	ChatID           string     `gorm:"column:chat_id;type:character varying(255);not null;default:NULL" json:"chat_id"`
	QuizCityName     *string    `gorm:"column:quiz_city_name;type:character varying(255)" json:"quiz_city_name"`
	DateQuizFinished *time.Time `gorm:"column:date_quiz_finished;type:timestamp without time zone" json:"date_quiz_finished"`
}

// var textWinner string = "Запуск протокола \"Победитель\"...\nChat_bot_assistant_260495\nАнализ профиля...\nСборка данных завершена.\n\nСнова приветствую тебя, друг! Фестиваль \"Действуй\" завершен, но твоя история еще продолжается. \n\nМоя система данных зафиксировала, что ты являешься первым собравшим все 5 QR-кодов, а следовательно открывшим свой путеводитель целиком.\n\nПрими мои поздравления, ведь ты стал(а) счастливым обладателем подарка😌\n\nТвой подарок будет в твоём населенном пункте до 30 апреля. Как только он поступит, моя система уведомит тебя об этом, и подскажет, где ты сможешь его забрать.\n\nПоздравляю и удачи!"
// var textLoser string = "Запуск протокола \"Возвращение\"...\nChat_bot_assistant_260495\nАнализ профиля...\nСборка данных завершена.\n\nСнова приветствую тебя, друг! Фестиваль \"Действуй\" завершен, но твоя история еще продолжается.\n\nМоя система проводила розыгрыш подарков для первых собравших свой путеводитель из пяти QR-кодов. \n\nК сожалению, в этот раз тебе не удалось стать победителем, но я очень ценю твое участие в Фестивале. Моя нейронная сеть очень надеется увидеть тебя на наших будущих мероприятиях в Югре!\n\nЯ желаю тебе удачи, ты - молодец!"
var textWinner string = `ПОВТОРНЫЙ ЗАПУСК протокола \"ПОБЕДИТЕЛЬ\"\.
Chat\_bot\_assistant\_260495
Анализ профиля\.\.\.
Сборка данных завершена\.

Снова приветствую тебя, друг\! Мы обещали подарки победителям и время пришло\!
Если ты еще не получил свой подарок, то переходи по ссылке, ищи адрес в своем городе и приходи за подарком 😌

https://telegra\.ph/FESTIVAL\-DEJSTVUJ\-05\-05

Поздравляю и удачи\!

А чтобы не упустить возможность окунуться в мир новых открытий и невероятных возможностей — не забудь подписаться на социальные сети Молодежного центра Югры:

— [ВКонтакте](https://vk\.com/mcugra)
— [Телеграм\-канал](https://t\.me/mc\_ugra)`

var ctx = context.Background()
var ncUser *nats.NatsClient = nats.Init(nats.NatsSettings{
	Ctx: ctx,
	URL: config.NATS_URL,
})

// var db = pg.DB

func main() {
	wg := &sync.WaitGroup{}

	ncUser.UsePublishSubject(config.NATS_USER_MESSAGES_SUBJECT)
	winnersFile, err := os.Open("./scripts/send_winners/winners.json")
	if err != nil {
		panic(err)
	}
	winnersBytes, _ := io.ReadAll(winnersFile)
	winners := []User{}

	json.Unmarshal(winnersBytes, &winners)

	fmt.Println(winners[0])

	// create slice of winnerIds
	var winnerIds []string
	for _, winner := range winners {
		winnerIds = append(winnerIds, winner.ChatID)
	}

	winnerIds = append(winnerIds, "306562182")
	winnerIds = append(winnerIds, "222414873")

	messageWinners := &models.SendableMessage{
		Text: &textWinner,
		SendOptions: &telebot.SendOptions{
			ParseMode: "MarkdownV2",
		},
	}

	// losers, err := GetLosers(winnerIds)
	// if err != nil {
	// 	panic(err)
	// }
	// loserIds := []string{}
	// // create slice of loserIds
	// for _, loser := range losers {
	// 	loserIds = append(loserIds, loser.ChatID)
	// }

	// loserIds = append(loserIds, "306562182")

	// messageLosers := &models.SendableMessage{
	// 	Text: &textLoser,
	// }
	wg.Add(2)
	sendAll(&winnerIds, messageWinners, ncUser)
	// sendAll(&loserIds, messageLosers, ncUser)
	wg.Wait()

}

func sendAll(targetChatIDs *[]string, message *models.SendableMessage, nc *nats.NatsClient) {
	var count int = 0

	for _, chatID := range *targetChatIDs {

		count += 1
		if count%30 == 0 {
			time.Sleep(1 * time.Second)
			fmt.Println("Count: ", count, "\n Sleeping...")
		}

		message.Recipient = &models.Recipient{
			ChatID: chatID,
		}
		nc.Publish(message)
	}
	fmt.Println("All sent!")
}

// func GetLosers(winnerIds []string) ([]User, error) {
// 	var users []User
// 	adminChatIds := db.GetAdminChatIDs()

// 	result := db.Model(&dbmodels.User{}).
// 		Select("users.chat_id").
// 		Where("users.chat_id NOT IN (?) AND users.chat_id NOT IN (?)", adminChatIds, winnerIds).
// 		Find(&users)

// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

// 	return users, nil
// }
