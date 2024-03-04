package userbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"quree/config"
	"quree/internal/models"
	"quree/internal/models/enums"
	"quree/internal/pg/dbmodels"
	"quree/internal/s3"
	"sort"
	"time"

	"quree/internal/pg"

	tele "gopkg.in/telebot.v3"

	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"
)

var db = pg.DB

func startHandler(c tele.Context) error {

	msgs := db.GetMessagesByType(enums.START)

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Sort < msgs[j].Sort
	})

	for _, m := range msgs {

		msg, err := json.Marshal(&models.MessageWithRecipient{
			ChatID:  fmt.Sprint(c.Chat().ID),
			Message: m,
		})
		if err != nil {
			log.Println(err)
			continue
		}

		nc.NC.Publish(config.NATS_MESSAGES_SUBJECT, msg)
	}

	return nil
}

func idHandler(c tele.Context) error {

	for i := 0; i < 20; i++ {

		msg, err := json.Marshal(&models.MessageWithRecipient{
			ChatID: fmt.Sprint(c.Chat().ID),
			Message: models.Message{
				// content is i with time.Now() with seconds and milliseconds
				Content: fmt.Sprint("ID: ", i, " Sent at: ", time.Now().Format("15:04:05.000")),
			},
		})
		if err != nil {
			log.Println(err)
			continue
		}

		nc.NC.Publish(config.NATS_MESSAGES_SUBJECT, msg)
	}

	return nil

	// return c.Send(fmt.Sprintf("%d", c.Chat().ID))
}

func qrHandler(c tele.Context) error {
	chatID := fmt.Sprint(c.Chat().ID)
	user := db.GetUserByChatIDAndRole(chatID, enums.USER)

	file := db.GetFileRecordByID(user.QRCode)
	sm := models.CreateSendableMessage(SendLimiter, &models.Message{
		Content: "Your QRCode",
	}, file)

	return sm.Send(c.Bot(), c.Chat(), &tele.SendOptions{})
}

func registerHandler(c tele.Context) error {

	chatID := fmt.Sprint(c.Chat().ID)

	user := db.GetUserByChatIDAndRole(chatID, enums.USER)

	if user != nil {
		sm := models.CreateSendableMessage(SendLimiter, &models.Message{
			Content: "Вы уже зарегистрированы! /start /qr /register",
		}, nil)

		return sm.Send(c.Bot(), c.Chat(), &tele.SendOptions{})
	}

	qrCodeUUID := uuid.New().String()
	qrCodeWidth := int32(256)
	qrCodeHeight := qrCodeWidth

	qrCodeMessage := &models.QRCodeMessage{
		ChatID:   chatID,
		QRCodeID: models.UUID(qrCodeUUID),
	}

	qrCodeMessageJSON, err := json.Marshal(qrCodeMessage)
	if err != nil {
		return err
	}

	fmt.Println(string(qrCodeMessageJSON))

	png, err := qrcode.Encode(string(qrCodeMessageJSON), qrcode.Medium, int(qrCodeWidth))
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
		return c.Send(err.Error())
	}

	err = db.CreateUser(&models.User{
		ChatID:      fmt.Sprint(c.Chat().ID),
		PhoneNumber: "test",
		Role:        enums.USER,
		QRCode:      models.UUID(qrCodeUUID),
	})

	if err != nil {
		return c.Send(err.Error())
	}

	sm := models.CreateSendableMessage(SendLimiter, &models.Message{
		Content: "Вы зарегистрированы! /start /qr /register",
	}, nil)

	return sm.Send(c.Bot(), c.Chat(), &tele.SendOptions{})
}

func CheckAuthorize() tele.MiddlewareFunc {
	l := log.Default()

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Message().Text == "/register" {
				return next(c)
			}

			chatID := fmt.Sprint(c.Chat().ID)
			user := db.GetUserByChatIDAndRole(chatID, enums.USER)

			if user == nil {
				sm := models.CreateSendableMessage(SendLimiter, &models.Message{
					Content: "Вы не зарегистрированы! Для регистрации введите /register",
				}, nil)

				l.Println("Юзер", chatID, "не авторизован")
				return sm.Send(c.Bot(), c.Chat(), &tele.SendOptions{})
			}

			l.Println("Юзер", chatID, "авторизован")
			return next(c)
		}
	}
}
