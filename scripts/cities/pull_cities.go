package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/skip2/go-qrcode"
)

type City struct {
	Name string `json:"name"`
}

func main() {

	// write a code that will load cities from json file cities.json, which located in ./cities.json
	// and print them to the console

	citiesFile, err := os.Open("./scripts/cities/cities.json")
	citiesBytes, _ := io.ReadAll(citiesFile)
	citiesJSON := []City{}

	if err != nil {
		panic(err)
	}

	defer citiesFile.Close()

	err = json.Unmarshal(citiesBytes, &citiesJSON)
	if err != nil {
		panic(err)
	}

	for _, city := range citiesJSON {
		fmt.Println(city)
		qrCodeWidth := int32(512)

		qrCodeMessage := fmt.Sprintf("loc,%s", city.Name)

		err := qrcode.WriteFile(
			qrCodeMessage,
			qrcode.Medium,
			int(qrCodeWidth),
			fmt.Sprintf("./scripts/cities/qrcodes/%s.png", city.Name),
		)

		if err != nil {
			panic(err)
		}
	}
}
