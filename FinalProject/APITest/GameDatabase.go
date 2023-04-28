package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	r := strings.NewReader("search \"zelda\"; fields *; limit 1;")
	req, _ := http.NewRequest("POST", "https://api.igdb.com/v4/games", r)
	client := &http.Client{}
	req.Header.Set("Client-ID", os.Getenv("twitchAPIID"))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("twitchAPIAccess")))

	resp, erro := client.Do(req)
	// fmt.Println(req)
	if erro != nil {
		log.Fatal(erro)
	}
	io.Copy(os.Stdout, resp.Body)
	// fmt.Println(resp.Body)
	defer resp.Body.Close()
}
