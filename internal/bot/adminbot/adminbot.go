package adminbot

import (
	"context"
	"fmt"
	"log"
	"quree/config"
	"time"

	"quree/internal/helpers"
	"quree/internal/models"
	"quree/internal/nats"
	"quree/internal/pg"
	"quree/internal/pg/dbmodels"
	"quree/internal/sendlimiter"

	"github.com/google/uuid"
	tele "gopkg.in/telebot.v3"
)

var ctx = context.Background()
var sl = sendlimiter.Init(ctx, config.RATE_LIMIT_GLOBAL, config.RATE_LIMIT_BURST_GLOBAL)
var db = pg.DB
var nc *nats.NatsClient = nats.Init(nats.NatsSettings{
	Ctx: ctx,
	URL: config.NATS_URL,
})

func Init() *tele.Bot {

	token := config.ADMIN_BOT_TOKEN

	log.Println("bot token:", token)

	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Use(helpers.RateLimit(sl))
	bot.Use(helpers.BotMiniLogger())
	bot.Use(CheckAuthorize())

	bot.Handle("/id", idHandler)
	bot.Handle("/start", registerHandler)
	bot.Handle(config.ADMIN_AUTH_CODE, registerHandler)

	// handle buttons
	bot.Handle(replyButtons["start"], registerHandler)
	bot.Handle(replyButtons["help"], helpHandlerFactory(nil))
	bot.Handle(inlineButtons["help1"], helpHandlerFactory(inlineButtons["help1"]))
	bot.Handle(inlineButtons["help2"], helpHandlerFactory(inlineButtons["help2"]))
	bot.Handle(inlineButtons["help3"], helpHandlerFactory(inlineButtons["help3"]))

	nc.UsePublishSubject(config.NATS_ADMIN_MESSAGES_SUBJECT)

	return bot
}

func idHandler(c tele.Context) error {
	return c.Send(fmt.Sprintf("%d", c.Chat().ID))
}

func helpHandlerFactory(btn *tele.InlineButton) tele.HandlerFunc {
	return func(c tele.Context) error {
		var chatID = fmt.Sprint(c.Chat().ID)

		var text *string
		var photo *models.Photo
		var inlineKeyboard [][]tele.InlineButton
		var data string
		if btn != nil {
			data = btn.Data
		} else {
			data = ""
		}

		switch data {
		case "help1":
			text = &textHelp1
			inlineKeyboard = [][]tele.InlineButton{{*inlineButtons["help2"]}}
		case "help2":
			text = nil
			photo = &models.Photo{File: tele.FromURL(passExampleURL), Caption: textHelp2}
			inlineKeyboard = [][]tele.InlineButton{{*inlineButtons["help3"]}}
		case "help3":
			text = &textHelp3
		default:
			text = nil
			photo = &models.Photo{File: tele.FromURL(scannerAdminExampleURL), Caption: textHelp}
			inlineKeyboard = [][]tele.InlineButton{{*inlineButtons["help1"]}}
		}

		var message = &models.SendableMessage{
			Text: text,
			Recipient: &models.Recipient{
				ChatID: chatID,
			},
			Photo: photo,
			SendOptions: &tele.SendOptions{
				ReplyMarkup: &tele.ReplyMarkup{
					InlineKeyboard: inlineKeyboard,
				},
			},
		}

		c.Respond()
		nc.Publish(message)

		return nil
	}
}

func registerHandler(c tele.Context) error {

	chatID := fmt.Sprint(c.Chat().ID)
	user := db.GetAdminByChatID(chatID)
	var text string

	if user != nil {
		text = textAlreadyRegistered
	} else {
		text = textRegistered
		err := db.CreateAdmin(&dbmodels.Admin{
			ID:          uuid.New().String(),
			DateCreated: time.Now(),
			ChatID:      chatID,
		})
		if err != nil {
			return err
		}
	}

	var message = &models.SendableMessage{
		Text: &text,
		Recipient: &models.Recipient{
			ChatID: chatID,
		},
		SendOptions: &tele.SendOptions{
			ReplyMarkup: menuAuthorized,
		},
	}

	nc.Publish(message)

	c.Bot().SetMenuButton(c.Sender(), tele.MenuButtonDefault)

	return nil
}

func CheckAuthorize() tele.MiddlewareFunc {
	l := log.Default()

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Message().Text == config.ADMIN_AUTH_CODE {
				return next(c)
			}

			chatID := fmt.Sprint(c.Chat().ID)
			user := db.GetAdminByChatID(chatID)

			if user == nil {
				var message = &models.SendableMessage{
					Text: &textUnauthorized,
					Recipient: &models.Recipient{
						ChatID: chatID,
					},
					SendOptions: &tele.SendOptions{
						ReplyMarkup: menuUnauthorized,
					},
				}

				nc.Publish(message)

				c.Bot().SetMenuButton(c.Sender(), tele.MenuButtonDefault)

				return nil
			}

			l.Println("Админ", chatID, "авторизован")
			return next(c)
		}
	}
}

var webAppScanner = &tele.WebApp{
	URL: config.ADMIN_WEBAPP_URL,
}

var webAppPassCheck = &tele.WebApp{
	URL: config.ADMIN_PASS_CHECK_WEBAPP_URL,
}

var textUnauthorized = "Вы не авторизованы! Для доступа к приложению введите код, полученный у куратора"
var textRegistered = "Вы авторизированы как админ! Нажмите Сканер QR для сканирования"
var textAlreadyRegistered = "Вы уже авторизированы! Нажмите Сканер QR для сканирования"

var scannerAdminExampleURL = "https://quree.maxnrm.monster/images/3835e4b8-d27e-4639-b493-868d93238ae7.png"
var textHelp = "Нажмите кнопку в боте \"Сканер QR\", откроется камера, попросите участника показать QR-код участника (у участника он открывается в его боте, по кнопке \"Показать QR-код\"), отсканируйте QR-код, участнику должно прийти сообщение.\n\nВнимание! Сканировать нужно только сканером по кнопке \"Сканер QR\" в боте! Обычная камера телефона не подойдет."
var textHelp1 = "На викторине участники сканируют QR-код викторины сами, из своего бота, просто помогите им найти QR-код викторины."

var passExampleURL = "https://quree.maxnrm.monster/images/7c813d8f-ed42-4921-a931-c715252e9ab2.png"
var textHelp2 = "На финальном мероприятии просканируйте QR-код участника с помощью кнопки \"Проверить проходку\", если ему пришла проходка, то он прошел все этапы, если нет, то нет."
var textHelp3 = "Отлично, если что, всегда можно перечитать еще раз!\n\nВ самом крайнем случае, если совершенно ничего непонятно или бот не работает, пишите сюда @gakhh в таком формате:\n- название города / поселения\n- точное описание проблемы\n- скриншот проблемы\n\nУдачи!"

var menuAuthorized = &tele.ReplyMarkup{
	ResizeKeyboard: true,
	ReplyKeyboard:  replyKeyboardAuthorized,
}

var menuUnauthorized = &tele.ReplyMarkup{
	ResizeKeyboard: true,
	ReplyKeyboard:  replyKeyboardUnauthorized,
}

var replyKeyboardAuthorized = [][]tele.ReplyButton{
	{*replyButtons["help"]},
	{*replyButtons["scanner"]},
	{*replyButtons["pass_check"]},
}

var replyKeyboardUnauthorized = [][]tele.ReplyButton{
	{*replyButtons["start"]},
}

var replyButtons = map[string]*tele.ReplyButton{
	"start":      {Text: "Начать"},
	"help":       {Text: "Как это работает?"},
	"scanner":    {Text: "Сканер QR", WebApp: webAppScanner},
	"pass_check": {Text: "Проверить проходку", WebApp: webAppPassCheck},
}

var inlineButtons = map[string]*tele.InlineButton{
	"help1": {Text: "А как с викториной?", Unique: "inlinehelp1", Data: "help1"},
	"help2": {Text: "А финальное мероприятие?", Unique: "inlinehelp2", Data: "help2"},
	"help3": {Text: "Вроде все понял :)", Unique: "inlinehelp3", Data: "help3"},
}
