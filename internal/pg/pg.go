package pg

import (
	"errors"
	"fmt"
	"quree/config"
	"time"

	"quree/internal/helpers"
	"quree/internal/models"
	"quree/internal/models/enums"

	"quree/internal/pg/dbmodels"
	"quree/internal/pg/dbquery"

	"github.com/google/uuid"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type pg struct {
	*gorm.DB
	q *dbquery.Query
}

var DB = Init(config.POSTGRES_CONN_STRING)

func Init(connString string) *pg {

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprint("Failed to connect to database at dsn: ", connString))
	}

	var pg *pg = &pg{
		DB: db,
		q:  dbquery.Use(db),
	}

	return pg
}

// function to create user
func (pg *pg) CreateUser(user *models.User) error {

	result := pg.Create(&dbmodels.User{
		ID:          uuid.New().String(),
		ChatID:      user.ChatID,
		PhoneNumber: &user.PhoneNumber,
		DateCreated: time.Now(),
		Role:        string(user.Role),
		QrCode:      string(user.QRCode),
	})

	return result.Error
}

// function to get user from db using ChatID, transform in models.User struct and return
func (pg *pg) GetUserByChatID(chatID string) *models.User {

	var user dbmodels.User
	result := pg.Where("chat_id = ?", chatID).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	var pn string
	if user.PhoneNumber != nil {
		pn = *user.PhoneNumber
	} else {
		pn = ""
	}

	return &models.User{
		ChatID:      user.ChatID,
		PhoneNumber: pn,
		Role:        helpers.StringToUserRole(user.Role),
		QRCode:      models.UUID(user.QrCode),
	}
}

// function to UploadFile in s3 and create record in db in table Files

func (pg *pg) CreateFileRecord(file *dbmodels.File) error {

	result := pg.Create(file)

	return result.Error
}

// get file by id
func (pg *pg) GetFileRecordByID(id models.UUID) *models.File {

	var file dbmodels.File
	result := pg.Where("id = ?", id).First(&file)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	var title, ftype string
	if file.Title != nil {
		title = *file.Title
	} else {
		title = ""
	}

	if file.Type != nil {
		ftype = *file.Type
	} else {
		ftype = ""
	}

	return &models.File{
		ID:       models.UUID(file.ID),
		Filename: *file.FilenameDisk,
		Title:    title,
		Type:     ftype,
	}
}

// function to get messages by type

func (pg *pg) GetMessagesByType(messageType enums.MessageType) []models.Message {

	var messages []dbmodels.Message
	result := pg.Where("type = ?", messageType).Find(&messages)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	var msgs []models.Message
	for _, message := range messages {

		var messageContent, messageImage string
		var messageSort int32

		if message.Content != nil {
			messageContent = *message.Content
		} else {
			messageContent = ""
		}

		if message.Image != nil {
			messageImage = *message.Image
		} else {
			messageImage = ""
		}

		if message.Sort != nil {
			messageSort = *message.Sort
		} else {
			messageSort = 0
		}

		msgs = append(msgs, models.Message{
			Content: messageContent,
			Image:   models.UUID(messageImage),
			Type:    messageType,
			Sort:    messageSort,
		})
	}

	return msgs
}

func (pg *pg) Close() {
	sqlDB, err := pg.DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()
}
