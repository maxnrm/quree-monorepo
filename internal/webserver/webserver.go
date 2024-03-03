// create basic gin server

package webserver

import (
	"encoding/json"
	"fmt"
	"quree/internal/models"
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

	var visit models.UserEventVisit

	err := c.BindJSON(&visit)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}

	data, _ := json.Marshal(visit)

	fmt.Println(string(data))

	db.CreateUserEventVisit(&visit)

	c.JSON(201, gin.H{"status": "created"})
}
