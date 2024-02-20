package bot

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"quree/config"
	"time"

	tele "gopkg.in/telebot.v3"
)

type CatFact struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

func Init() *tele.Bot {

	log.Println("bot token:", config.ADMIN_BOT_TOKEN)

	settings := tele.Settings{
		Token:  config.ADMIN_BOT_TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(settings)
	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/hello", helloHandler)
	b.Handle("кот", quoteHandler)

	return b
}

func helloHandler(c tele.Context) error {
	return c.Send("Hello!")
}

func quoteHandler(c tele.Context) error {
	url := "https://catfact.ninja/fact"

	req, _ := http.NewRequest("GET", url, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Error while getting quote: ", err)
	}

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	var catFact CatFact

	err = json.Unmarshal(body, &catFact)
	if err != nil {
		log.Fatal("Error while unmarshalling json: ", err)
	}

	return c.Send(catFact.Fact)
}
