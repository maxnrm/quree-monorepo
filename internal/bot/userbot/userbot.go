package userbot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"quree/config"
	"quree/internal/helpers"
	"quree/internal/models"
	"quree/internal/models/enums"
	"quree/internal/nats"
	"quree/internal/pg"
	"quree/internal/pg/dbmodels"
	"quree/internal/s3"
	"quree/internal/sendlimiter"
	"time"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	tele "gopkg.in/telebot.v3"
)

var ctx = context.Background()
var SendLimiter = sendlimiter.Init(ctx)
var db = pg.DB
var nc *nats.NatsClient = nats.Init(nats.NatsSettings{
	Ctx: ctx,
	URL: config.NATS_URL,
})

func Init() *tele.Bot {

	token := config.USER_BOT_TOKEN

	log.Println("bot token:", token)

	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Use(helpers.BotMiniLogger())
	bot.Use(CheckAuthorize())

	bot.Handle("/start", startHandler)
	bot.Handle("/help", helpHandler)
	bot.Handle("/id", idHandler)
	bot.Handle("/qr", qrHandler)
	bot.Handle("/register", registerHandler)

	nc.UsePublishSubject(config.NATS_USER_MESSAGES_SUBJECT)

	return bot
}

func startHandler(c tele.Context) error {

	chatID := fmt.Sprint(c.Chat().ID)

	rec := &models.Recipient{
		ChatID: chatID,
	}

	messages := db.GetMessagesByType(enums.START)
	message := messages[0]
	message.Recipient = rec
	json, err := json.Marshal(message)
	if err != nil {
		return err
	}

	nc.Publish(json)

	return nil
}

func helpHandler(c tele.Context) error {

	chatID := fmt.Sprint(c.Chat().ID)

	rec := &models.Recipient{
		ChatID: chatID,
	}

	messages := db.GetMessagesByType(enums.START)
	message := messages[0]
	message.Recipient = rec
	json, err := json.Marshal(message)
	if err != nil {
		return err
	}

	nc.Publish(json)

	return nil
}

func idHandler(c tele.Context) error {
	chatID := fmt.Sprint(c.Chat().ID)
	rec := &models.Recipient{
		ChatID: chatID,
	}

	for i := 0; i < 20; i++ {

		text := fmt.Sprintf("index: %d", i)

		var message = models.SendableMessage{
			Text: &text,
		}
		message.Recipient = rec

		json, err := json.Marshal(message)
		if err != nil {
			return err
		}

		nc.Publish(json)
	}

	return nil
}

func qrHandler(c tele.Context) error {
	chatID := fmt.Sprint(c.Chat().ID)
	user := db.GetUserByChatID(chatID)
	if user == nil {
		return errors.New("user not found")
	}

	photoURL := config.IMGPROXY_PUBLIC_URL + "/" + user.QrCode + ".png"

	qr := &tele.Photo{File: tele.FromURL(photoURL), Caption: "Твой QR-код"}

	message := models.SendableMessage{Photo: qr}

	rec := &models.Recipient{
		ChatID: chatID,
	}

	message.Recipient = rec

	json, err := json.Marshal(message)
	if err != nil {
		return err
	}

	nc.Publish(json)

	return nil
}

func registerHandler(c tele.Context) error {

	chatID := fmt.Sprint(c.Chat().ID)

	user := db.GetUserByChatID(chatID)
	if user != nil {
		return helpHandler(c)
	}

	// default middleware should not return anything, if user is not registered

	qrCodeUUID := uuid.New().String()
	qrCodeWidth := int32(256)
	qrCodeHeight := qrCodeWidth

	qrCodeMessage := fmt.Sprintf("%s,%s", chatID, qrCodeUUID)

	png, err := qrcode.Encode(qrCodeMessage, qrcode.Medium, int(qrCodeWidth))
	if err != nil {
		return err
	}

	qrCodeBytesReader := bytes.NewReader(png)
	qrCodeSize := qrCodeBytesReader.Size()

	filenameDisk := fmt.Sprintf("%s.png", qrCodeUUID)
	filenameDownload := fmt.Sprintf("%s.png", chatID)
	fileType := "image/png"

	info, err := s3.S3Client.UploadImage(filenameDisk, qrCodeBytesReader, qrCodeSize)
	if err != nil {
		return err
	}

	err = db.CreateFileRecord(&dbmodels.File{
		Storage:          "s3",
		ID:               qrCodeUUID,
		Title:            &chatID,
		FilenameDisk:     &filenameDisk,
		FilenameDownload: filenameDownload,
		UploadedOn:       info.LastModified,
		Filesize:         &qrCodeSize,
		Width:            &qrCodeWidth,
		Height:           &qrCodeHeight,
		Type:             &fileType,
	})

	if err != nil {
		fmt.Println("error creating qr:", err)
	}

	err = db.CreateUser(&dbmodels.User{
		ID:          uuid.New().String(),
		DateCreated: time.Now(),
		ChatID:      chatID,
		QrCode:      qrCodeUUID,
	})
	if err != nil {
		fmt.Println("error creating user:", err)
	}

	return helpHandler(c)
}

func CheckAuthorize() tele.MiddlewareFunc {
	l := log.Default()
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Message().Text == "/register" {
				return next(c)
			}

			chatID := fmt.Sprint(c.Chat().ID)
			user := db.GetUserByChatID(chatID)

			if user == nil {
				return startHandler(c)
			}

			l.Println("Юзер", chatID, "авторизован")
			return next(c)
		}
	}
}
