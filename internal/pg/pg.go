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

	var result *gorm.DB

	if user.Role == enums.ADMIN {
		result = pg.Create(&dbmodels.User{
			ID:          uuid.New().String(),
			ChatID:      &user.ChatID,
			PhoneNumber: &user.PhoneNumber,
			DateCreated: time.Now(),
			Role:        string(user.Role),
			QrCode:      nil,
		})

	} else if user.Role == enums.USER {
		qrCodeStr := string(user.QRCode)

		result = pg.Create(&dbmodels.User{
			ID:          uuid.New().String(),
			ChatID:      &user.ChatID,
			PhoneNumber: &user.PhoneNumber,
			DateCreated: time.Now(),
			Role:        string(user.Role),
			QrCode:      &qrCodeStr,
		})
	}

	return result.Error

}

// function to get user from db using ChatID, transform in models.User struct and return
func (pg *pg) GetUserByChatIDAndRole(chatID string, role enums.UserRole) *models.User {

	var user dbmodels.User
	// get user by chatID and Role
	result := pg.Where("chat_id = ? AND role = ?", chatID, string(role)).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	var userQRCode string
	if user.QrCode == nil {
		userQRCode = ""
	} else {
		userQRCode = *user.QrCode
	}

	var phoneNumber string
	if user.PhoneNumber == nil {
		phoneNumber = ""
	} else {
		phoneNumber = *user.PhoneNumber
	}

	return &models.User{
		ChatID:      *user.ChatID,
		PhoneNumber: phoneNumber,
		Role:        enums.UserRole(user.Role),
		QRCode:      models.UUID(userQRCode),
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

func (pg *pg) CreateUserEventVisit(visit models.UserEventVisit) error {

	visitAdminID := string(visit.AdminID)
	visitQuizID := string(visit.QuizID)

	result := pg.Create(&dbmodels.UserEventVisit{
		ID:          uuid.New().String(),
		UserID:      string(visit.UserID),
		DateCreated: time.Now(),
		EventType:   string(visit.Type),
		AdminID:     &visitAdminID,
		QuizID:      &visitQuizID,
	})

	return result.Error
}

// method to count UserEventVisits for user, counting only events with type EVENT

func (pg *pg) CountUserEventVisitsForUser(userID models.UUID) int64 {

	var count int64
	result := pg.Model(&dbmodels.UserEventVisit{}).Where("user_id = ? AND event_type = ?", userID, string(enums.EVENT)).Count(&count)

	if result.Error != nil {
		return 0
	}

	return count
}

func (pg *pg) Close() {
	sqlDB, err := pg.DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()
}
