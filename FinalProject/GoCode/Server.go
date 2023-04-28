package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	googleapi "google.golang.org/api/googleapi"
	option "google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

func main() {
	var sd serverData
	sd.TwitchID = os.Getenv("twitchAPIID")
	sd.TwitchAccess = os.Getenv("twitchAPIAccess")
	sd.YoutubeKey = os.Getenv("youtubeAPIKey")

	fs := http.FileServer(http.Dir("../Webpage/"))
	http.Handle("/", fs)
	http.HandleFunc("/search", sd.search)
	log.Fatal(http.ListenAndServe(":8002", nil)) // Begin and log the server
}

type serverData struct {
	TwitchID     string
	TwitchAccess string
	YoutubeKey   string
}

type GameNameInput struct {
	Name string
}

type ResponseOutput struct {
	Result    string
	Name      string
	Storyline string
	Boxart    string
	VideoId   []string
}

type GameArray struct {
	Id        int
	Cover     int
	Name      string
	Storyline string
}

type CoverResponse struct {
	ID  int
	URL string
}

func (sd *serverData) search(w http.ResponseWriter, req *http.Request) {
	var gameName GameNameInput
	if req.Method == "POST" {
		// Read the request
		if err := json.NewDecoder(req.Body).Decode(&gameName); err != nil {
			log.Println(req.Body)
			var response ResponseOutput
			response.Result = "nok"
			response.Name = "void"
			response.VideoId = nil
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
			log.Fatal(err.Error())
		}
		log.Println("Requested Name = " + gameName.Name)

		// Get the closest game with that name
		r := strings.NewReader(fmt.Sprintf("search \"%s\"; fields name,cover,storyline; limit 1;", gameName.Name))
		gamereq, _ := http.NewRequest("POST", "https://api.igdb.com/v4/games", r)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		gamereq.Header.Set("Client-ID", os.Getenv("twitchAPIID"))
		gamereq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("twitchAPIAccess")))

		gameresp, erro := client.Do(gamereq)
		// fmt.Println(req)
		if erro != nil {
			log.Fatal(erro.Error())
		}
		defer gameresp.Body.Close()

		var gamearray []GameArray
		if err := json.NewDecoder(gameresp.Body).Decode(&gamearray); err != nil {
			log.Fatal(err.Error())
		}

		log.Println(gamearray[0].Cover)
		r = strings.NewReader(fmt.Sprintf("fields url; where id=%d;", gamearray[0].Cover))
		gamereq, _ = http.NewRequest("POST", "https://api.igdb.com/v4/covers", r)
		gamereq.Header.Set("Client-ID", os.Getenv("twitchAPIID"))
		gamereq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("twitchAPIAccess")))

		coverresp, erro := client.Do(gamereq)
		if erro != nil {
			log.Fatal(erro.Error())
		}
		defer coverresp.Body.Close()

		var cover []CoverResponse
		if err := json.NewDecoder(coverresp.Body).Decode(&cover); err != nil {
			log.Fatal(err.Error())
		}

		// Get the first gameplay video with that name
		developerKey := os.Getenv("youtubeAPIKey")

		ctx := context.Background()
		service, err := youtube.NewService(ctx, option.WithAPIKey(developerKey), option.WithHTTPClient(client))
		if err != nil {
			log.Fatalf("Error creating new YouTube client: %v", err)
		}

		// Make the API call to YouTube.
		Test := []string{"id", "snippet"}
		call := service.Search.List(Test).
			Q(fmt.Sprintf("%s gameplay", gamearray[0].Name)).
			MaxResults(10)
		ytresponse, err := call.Do(googleapi.QueryParameter("key", developerKey))
		if err != nil {
			log.Fatal(err.Error())
		}

		// Send response
		var response ResponseOutput
		response.Result = "ok"
		response.Name = gamearray[0].Name
		response.Storyline = gamearray[0].Storyline
		response.Boxart = cover[0].URL

		log.Println("Found Name = " + gamearray[0].Name)
		log.Println("Storyline = " + gamearray[0].Storyline)
		log.Println("Cover URL = " + cover[0].URL)

		for _, item := range ytresponse.Items {
			switch item.Id.Kind {
			case "youtube#video":
				response.VideoId = append(response.VideoId, item.Id.VideoId)
			}
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(&response); err != nil {
			log.Fatal(err.Error())
		}
	}
	defer req.Body.Close()
}
