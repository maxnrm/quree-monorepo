package userbot

import (
	"bytes"
	"fmt"
	"quree/internal/models"
	"quree/internal/models/enums"
	"quree/internal/pg/dbmodels"
	"quree/internal/s3"
	"sort"

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
		var sm *models.SendableMessage

		if m.Image != "" {
			file := db.GetFileRecordByID(m.Image)
			sm = models.CreateSendableMessage(SendLimiter, &m, file)
		} else {
			sm = models.CreateSendableMessage(SendLimiter, &m, nil)
		}

		err := sm.Send(c.Bot(), c.Chat(), &tele.SendOptions{})
		if err != nil {
			return err
		}

	}

	return nil
}

func idHandler(c tele.Context) error {
	return c.Send(fmt.Sprintf("%d", c.Chat().ID))
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
			Content: "Вы уже зарегистрированы!",
		}, nil)

		return sm.Send(c.Bot(), c.Chat(), &tele.SendOptions{})
	}

	qrCodeUUID := uuid.New().String()
	qrCodeWidth := int32(256)
	qrCodeHeight := qrCodeWidth

	png, err := qrcode.Encode(fmt.Sprintf("%s %s", chatID, qrCodeUUID), qrcode.Medium, int(qrCodeWidth))
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
		Content: "Registered!",
	}, nil)

	return sm.Send(c.Bot(), c.Chat(), &tele.SendOptions{})
}
