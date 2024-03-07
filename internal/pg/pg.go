package pg

import (
	"errors"
	"fmt"
	"quree/config"
	"time"

	"quree/internal/models"
	"quree/internal/models/enums"

	"quree/internal/pg/dbmodels"
	"quree/internal/pg/dbquery"

	tele "gopkg.in/telebot.v3"
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

func (pg *pg) CreateAdmin(admin *dbmodels.Admin) error {

	result := pg.Create(admin)

	return result.Error
}

func (pg *pg) CreateUser(user *dbmodels.User) error {

	result := pg.Create(user)

	return result.Error
}

// method, that adds city to user

func (pg *pg) UpdateUserQuizCity(userCityMessage *models.UserCityMessage) error {

	// get user by chatID
	user := pg.GetUserByChatID(userCityMessage.ChatID)
	if user == nil {
		return errors.New("user not found updateQuizCity")
	}

	// update user city
	user.QuizCityName = &userCityMessage.City
	user.IsFinished = true
	now := time.Now()
	user.DateQuizFinished = &now

	result := pg.Save(user)

	return result.Error
}

func (pg *pg) GetAdminByChatID(chatID string) *dbmodels.Admin {
	var admin *dbmodels.Admin
	// get user by chatID and Role
	result := pg.Where("chat_id = ?", chatID).First(&admin)
	if result.Error != nil {
		return nil
	}

	return admin
}

func (pg *pg) GetUserByChatID(chatID string) *dbmodels.User {
	var user *dbmodels.User
	// get user by chatID and Role
	result := pg.Where("chat_id = ?", chatID).First(&user)
	if result.Error != nil {
		return nil
	}

	return user
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
func (pg *pg) GetMessagesByType(messageType enums.MessageType) []*models.SendableMessage {

	var messages []*dbmodels.Message
	result := pg.Where("type = ?", messageType).Find(&messages)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	var sendableMessages []*models.SendableMessage

	for _, message := range messages {
		photo := &models.Photo{}
		if message.Image != nil {
			photoURL := config.IMGPROXY_PUBLIC_URL + "/" + *message.Image + ".jpg"
			photo = &models.Photo{File: tele.FromURL(photoURL)}
		}

		sendableMessages = append(sendableMessages, &models.SendableMessage{
			Text:  message.Text,
			Photo: photo,
		})
	}

	return sendableMessages
}

func (pg *pg) CreateUserEventVisit(visit *dbmodels.UserEventVisit) error {

	result := pg.Create(visit)

	return result.Error
}

// method to count UserEventVisits for user, counting only events with type EVENT

func (pg *pg) CountUserEventVisitsForUser(userChatID string) int64 {

	var count int64
	result := pg.Model(&dbmodels.UserEventVisit{}).Where("user_chat_id = ?", userChatID).Count(&count)

	if result.Error != nil {
		return 0
	}

	return count
}

// method to get latest UserEventVisit for user, counting only events with type EVENT

func (pg *pg) GetLatestUserEventVisitByUserChatID(userChatID string) (time.Time, error) {

	var visit dbmodels.UserEventVisit
	result := pg.Where("user_chat_id = ?", userChatID).Order("date_created desc").First(&visit)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return time.Now().Add(-10 * time.Minute), nil
	}

	return visit.DateCreated, nil
}

func (pg *pg) Close() {
	sqlDB, err := pg.DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()
}
