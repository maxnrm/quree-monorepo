// create basic gin server

package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"quree/config"
	"quree/internal/models"
	"quree/internal/models/enums"
	"quree/internal/nats"
	"quree/internal/pg"
	"time"

	"github.com/gin-gonic/gin"
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

	router.POST("/api/user_event_visit/create", CreateUserEventVisit)

	router.Run(":3000")

}

// create handler to add user event visit

func CreateUserEventVisit(c *gin.Context) {

	// json consists of two fields - user_id and admin_id

	var qrCodeMessage models.QRCodeMessage

	err := c.BindJSON(&qrCodeMessage)
	if err != nil {
		fmt.Println("Error binding json")
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}

	fmt.Println(qrCodeMessage)

	userID, err := db.GetUserIDByChatIDAndRole(qrCodeMessage.ChatID, enums.USER)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}

	latestEventVisit, _ := db.GetLatestUserEventVisitByUserID(userID)
	if time.Since(latestEventVisit).Minutes() < 5 {
		c.JSON(201, gin.H{"status": "scanned recently"})
		return
	}

	adminID, err := db.GetUserIDByChatIDAndRole(qrCodeMessage.AdminChatID, enums.ADMIN)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}

	visit := models.UserEventVisit{
		UserID:  userID,
		AdminID: adminID,
		Type:    enums.EVENT,
	}

	data, _ := json.Marshal(visit)

	fmt.Println(string(data))

	db.CreateUserEventVisit(&visit)

	numberOfEvents := db.CountUserEventVisitsForUser(userID)
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
		eventType = enums.LORE_EVENT_EXTRA
	}

	msgs := db.GetMessagesByType(eventType)

	if msgs != nil {
		fmt.Println("Messages not nil")
		//create a slice of msgs[i].Group to send to nats
		groupMap := make(map[int]bool)
		var uniqueGroups []int

		for _, obj := range msgs {
			if _, exists := groupMap[obj.Group]; !exists {
				groupMap[obj.Group] = true
				uniqueGroups = append(uniqueGroups, obj.Group)
			}
		}

		// Pick a random value from the uniqueGroups slice
		randomGroup := uniqueGroups[rand.Intn(len(uniqueGroups))]
		fmt.Println(uniqueGroups)
		var filtered []models.Message
		for _, obj := range msgs {
			if obj.Group == randomGroup {
				filtered = append(filtered, obj)
			}
		}

		for _, m := range filtered {

			fmt.Println(m.Sort)

			msg, err := json.Marshal(&models.MessageWithRecipient{
				ChatID:  qrCodeMessage.ChatID,
				Message: m,
			})
			if err != nil {
				log.Println(err)
				continue
			}

			nc.NC.Publish(config.NATS_MESSAGES_SUBJECT, msg)
		}
	}

	c.JSON(200, gin.H{"status": "created"})
}
