package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var tpl *template.Template

var (
	query      = flag.String("query", "Google", "Search term") // Parameters to change
	maxResults = flag.Int64("max-results", 5, "Max YouTube results")
)

var developerKey string

func main() {
	tpl, _ = tpl.ParseGlob("templates/responsePage.html")
	http.HandleFunc("/test", youtubeAPI)
	http.ListenAndServe(":8000", nil)
}

func youtubeAPI(w http.ResponseWriter, r *http.Request){
	searchList := []string{"id", "snippet"}
	flag.Parse()
	developerKey = os.Getenv("youtubeAPIKey")

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	call := service.Search.List(searchList).
		Q(*query).
		MaxResults(*maxResults)
	response, err := call.Do()
	//handleError(err, "")

	// Group video, channel, and playlist results in separate lists.
	videos := make(map[string]string)
	channels := make(map[string]string)
	playlists := make(map[string]string)

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		case "youtube#channel":
			channels[item.Id.ChannelId] = item.Snippet.Title
		case "youtube#playlist":
			playlists[item.Id.PlaylistId] = item.Snippet.Title
		}
	}

	printIDs("Videos", videos)
	printIDs("Channels", channels)
	printIDs("Playlists", playlists)

	tpl.ExecuteTemplate(w, "responsePage.html", response)
}

// Print the ID and title of each result in a list as well as a name that
// identifies the list. For example, print the word section name "Videos"
// above a list of video search results, followed by the video ID and title
// of each matching video.
func printIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for id, title := range matches {
		fmt.Printf("[%v] %v\n", id, title)
	}
	fmt.Printf("\n\n")
}
