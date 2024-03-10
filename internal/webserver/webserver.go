// create basic gin server

package webserver

import (
	"context"
	"fmt"
	"math/rand"
	"quree/config"
	"quree/internal/helpers"
	"quree/internal/models"
	"quree/internal/models/enums"
	"quree/internal/nats"
	"quree/internal/pg"
	"quree/internal/pg/dbmodels"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ctx = context.Background()
var db = pg.DB
var nc *nats.NatsClient = nats.Init(nats.NatsSettings{
	Ctx: ctx,
	URL: config.NATS_URL,
})

func Start() {

	fmt.Println("Starting webserver...")

	router := gin.Default()

	router.POST("/api/user_event_visit/create", createUserEventVisit)
	router.POST("/api/user/add_city", addUserCity)
	router.GET("/healthcheck", healthcheck)

	nc.UsePublishSubject(config.NATS_USER_MESSAGES_SUBJECT)

	router.Run(fmt.Sprintf(":%s", config.USER_WEBSERVER_PORT))

}

func healthcheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func addUserCity(c *gin.Context) {
	var userCityMessage models.UserCityMessage

	err := c.BindJSON(&userCityMessage)
	if err != nil {
		c.JSON(400, gin.H{"error": "Кажется, вы сканируете неправильный QR-код"})
		fmt.Println("error binding json userCityMessage")
		return
	}

	var chatID = userCityMessage.ChatID

	err = db.UpdateUserQuizCity(&userCityMessage)
	if err != nil {
		c.JSON(201, gin.H{"status": "Кажется, вы уже прошли викторину"})
		return
	}

	messages := db.GetMessagesByType(enums.LORE_EVENT_QUIZ)

	message := messages[0]
	message.Recipient = &models.Recipient{
		ChatID: chatID,
	}

	nc.Publish(message)

	c.JSON(200, gin.H{"status": "Вы успешно прошли викторину!"})
}

func createUserEventVisit(c *gin.Context) {

	// json consist of one significant field user_chat_id

	var qrCodeMessage models.QRCodeMessage

	err := c.BindJSON(&qrCodeMessage)
	if err != nil {
		fmt.Println("error binding json")
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}

	var today = time.Now()
	goalDate, err := time.Parse("2006-01-02", config.FINISH_PASS_DATE)
	if err != nil {
		fmt.Println("error parsing date")
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	if helpers.DateEqual(today, goalDate) {
		user := db.GetUserByChatID(qrCodeMessage.UserChatID)
		numberOfEvents := db.CountUserEventVisitsForUser(qrCodeMessage.UserChatID)

		if numberOfEvents > 4 && user.QuizCityName != nil {
			messages := db.GetMessagesByType(enums.FINAL_PASS)
			message := messages[0]
			message.Recipient = &models.Recipient{
				ChatID: qrCodeMessage.UserChatID,
			}

			nc.Publish(message)

			c.JSON(200, gin.H{"status": "accepted"})

			return
		}

	}

	userChatID := qrCodeMessage.UserChatID

	fmt.Println(qrCodeMessage)

	latestEventVisit, _ := db.GetLatestUserEventVisitByUserChatID(userChatID)
	if time.Since(latestEventVisit).Minutes() < float64(config.EVENT_VISIT_DELAY_MINUTES) {
		c.JSON(304, gin.H{"status": "scanned recently"})
		return
	}

	visit := dbmodels.UserEventVisit{
		ID:          uuid.New().String(),
		DateCreated: time.Now(),
		UserChatID:  userChatID,
	}

	db.CreateUserEventVisit(&visit)

	numberOfEvents := db.CountUserEventVisitsForUser(userChatID)
	var eventType enums.MessageType

	// switch number of events to send different messages

	switch {
	case numberOfEvents == 1:
		eventType = enums.LORE_EVENT1
	case numberOfEvents == 2:
		eventType = enums.LORE_EVENT2
	case numberOfEvents == 3:
		eventType = enums.LORE_EVENT3
	case numberOfEvents == 4:
		eventType = enums.LORE_EVENT4
	case numberOfEvents > 4:
		eventType = enums.LORE_EVENT_EXTRA
	default:
		c.JSON(500, gin.H{"status": "Что-то пошло не так..."})
		return
	}

	timeDBMessageStart := time.Now()
	messages := db.GetMessagesByType(eventType)
	timeDBMessage := time.Since(timeDBMessageStart)
	fmt.Println("Get message(s) in:", timeDBMessage.Milliseconds(), "ms")

	message := messages[rand.Intn(len(messages))]
	message.Recipient = &models.Recipient{
		ChatID: userChatID,
	}

	nc.Publish(message)

	c.JSON(201, gin.H{"status": "created"})
}
