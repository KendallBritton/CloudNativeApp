package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Access_token string
	Expires_in   int64
	Token_type   string
}

func main() {
	authUrl := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials", os.Getenv("twitchAPIID"), os.Getenv("twitchAPISecret"))
	req, _ := http.NewRequest("POST", authUrl, nil)
	client := &http.Client{}

	resp, erro := client.Do(req)
	fmt.Println(req)
	if erro != nil {
		log.Fatal(erro)
	}
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
}
