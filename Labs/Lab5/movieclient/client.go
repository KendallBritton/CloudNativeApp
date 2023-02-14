// Package main imlements a client for movieinfo service
package main

import (
	"context"
	"log"
	"os"
	"time"

	"labs/Lab5/movieapi"

	"google.golang.org/grpc"
)

const (
	address      = "localhost:50051"
	defaultTitle = "Pulp fiction"
)

type MovieData struct {
	title    string
	year     int32
	director string
	cast     []string
}

func main() {

	var newMovie MovieData

	newMovie.title = "The Dark Knight"
	newMovie.director = "Christopher Nolan"
	newMovie.year = 2008
	newMovie.cast = []string{"Christian Bale, Heath Ledger, Aaron Eckhart"}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := movieapi.NewMovieInfoClient(conn)

	// Contact the server and print out its response.
	title := defaultTitle
	if len(os.Args) > 1 {
		title = os.Args[1]
	}
	// Timeout if server doesn't respond
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: title})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("Movie Info for %s %d %s %v", title, r.GetYear(), r.GetDirector(), r.GetCast())
	test, err := c.SetMovieInfo(ctx, &movieapi.MovieData{Title: newMovie.title, Year: newMovie.year, Director: newMovie.director, Cast: newMovie.cast})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("%v", test.Code)
	testTitle := "The Dark Knight"
	testOutput, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: testTitle})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("Movie Info for %s %d %s %v", testTitle, testOutput.GetYear(), testOutput.GetDirector(), testOutput.GetCast())
	testTitle2 := "Fast and Furious"
	testOutput2, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: testTitle})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("Movie Info for %s %d %s %v", testTitle2, testOutput2.GetYear(), testOutput2.GetDirector(), testOutput2.GetCast())
}
