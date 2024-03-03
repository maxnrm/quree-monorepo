// create basic gin server

package webserver

import (
	"encoding/json"
	"fmt"
	"quree/internal/models"
	"quree/internal/models/enums"
	"quree/internal/pg"

	"github.com/gin-gonic/gin"
)

var db = pg.DB

func Start() {

	fmt.Println("Starting webserver...")

	router := gin.Default()

	router.POST("/api/user_event_visit/create", CreateUserEventVisit)

	router.Run(":8080")

}

// create handler to add user event visit

func CreateUserEventVisit(c *gin.Context) {

	// json consists of two fields - user_id and admin_id

	var qrCodeMessage models.QRCodeMessage

	err := c.BindJSON(&qrCodeMessage)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}

	userID, err := db.GetUserIDByChatIDAndRole(qrCodeMessage.ChatID, enums.USER)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}

	adminID, err := db.GetUserIDByChatIDAndRole(qrCodeMessage.ChatID, enums.ADMIN)
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

	c.JSON(201, gin.H{"status": "created"})
}
