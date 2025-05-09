package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"quree/internal/pg"
	"quree/internal/pg/dbmodels"
	"time"
)

var db = pg.DB

type winnerPerCity struct {
	CityName     string `json:"city_name"`
	WinnersCount int    `json:"winners_count"`
}

type User struct {
	ChatID           string     `gorm:"column:chat_id;type:character varying(255);not null;default:NULL" json:"chat_id"`
	QuizCityName     *string    `gorm:"column:quiz_city_name;type:character varying(255)" json:"quiz_city_name"`
	DateQuizFinished *time.Time `gorm:"column:date_quiz_finished;type:timestamp without time zone" json:"date_quiz_finished"`
}

func main() {
	winnersPerCityFile, err := os.Open("./scripts/get_winners/winners_by_city.json")
	if err != nil {
		panic(err)
	}
	winnerPerCityBytes, _ := io.ReadAll(winnersPerCityFile)
	winnersPerCity := []winnerPerCity{}
	var users []User
	adminChatIds := db.GetAdminChatIDs()

	json.Unmarshal(winnerPerCityBytes, &winnersPerCity)

	var count int = 0

	for _, c := range winnersPerCity {
		temp, err := GetWinnersByCity(c.CityName, c.WinnersCount, adminChatIds)
		if err != nil {
			panic(err)
		}

		fmt.Println("City:", c.CityName, "Should be:", c.WinnersCount, "Fact:", len(temp))

		count += c.WinnersCount

		users = append(users, temp...)
	}

	fmt.Println(count)
	fmt.Println(len(users))
	fmt.Println("City: ", *users[0].QuizCityName)

	err = SaveWinnersToFile(users, "winners.json")
	if err != nil {
		panic(err)
	}

}

func GetWinnersByCity(quizCityName string, limit int, adminChatIds []string) ([]User, error) {
	var users []User

	result := db.Model(&dbmodels.User{}).
		Select("users.chat_id, users.quiz_city_name, users.date_quiz_finished, COUNT(user_event_visits.id) AS VisitsCount").
		Joins("LEFT JOIN User_Event_Visits ON users.chat_id = user_event_visits.user_chat_id").
		Where("users.quiz_city_name = ? AND users.date_quiz_finished >= '2024-03-15'::date AND users.chat_id NOT IN (?)", quizCityName, adminChatIds).
		Group("users.chat_id, users.quiz_city_name, users.date_quiz_finished").
		Having("COUNT(user_event_visits.id) >= 4").
		Order("users.quiz_city_name DESC, users.date_quiz_finished ASC").
		Limit(limit).
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func SaveWinnersToFile(winners []User, filename string) error {
	winnersJson, err := json.Marshal(winners)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(winnersJson)
	if err != nil {
		return err
	}

	return nil
}
