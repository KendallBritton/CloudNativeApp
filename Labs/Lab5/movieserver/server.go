// Package main implements a server for movieinfo service.
package main

import (
	"context"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"

	"labs/Lab5/movieapi"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement movieapi.MovieInfoServer
type server struct {
	movieapi.UnimplementedMovieInfoServer
}

// Map representing a database
var moviedb = map[string][]string{"Pulp fiction": []string{"1994", "Quentin Tarantino", "John Travolta,Samuel Jackson,Uma Thurman,Bruce Willis"}}

// Add comment
func (s *server) SetMovieInfo(ctx context.Context, in *movieapi.MovieData) (*movieapi.Status, error) {
	title := in.Title
	year := in.Year
	cast := in.Cast
	director := in.Director
	status := &movieapi.Status{}

	if title == "" {
		return status, errors.New("Error in saving title")
	} else {
		status.Code = "Passed Successfully"
	}

	if director == "" {
		return status, errors.New("Error in saving director")
	} else {
		status.Code = "Passed Successfully"
	}

	for i := range cast {

		if cast[i] == "" {
			return status, errors.New("Error in saving cast")
		} else {
			status.Code = "Passed Successfully"
		}

	}

	stringYear := strconv.Itoa(int(year))

	temp := []string{stringYear, director}

	moviedb[title] = append(temp, cast...)

	log.Printf("Movie Data Saved: %v", title)

	return status, nil
}

// GetMovieInfo implements movieapi.MovieInfoServer
func (s *server) GetMovieInfo(ctx context.Context, in *movieapi.MovieRequest) (*movieapi.MovieReply, error) {
	title := in.GetTitle()
	log.Printf("Received: %v", title)
	reply := &movieapi.MovieReply{}
	if val, ok := moviedb[title]; !ok { // Title not present in database
		return reply, nil
	} else {
		if year, err := strconv.Atoi(val[0]); err != nil {
			reply.Year = -1
		} else {
			reply.Year = int32(year)
		}
		reply.Director = val[1]
		cast := strings.Split(val[2], ",")
		reply.Cast = append(reply.Cast, cast...)

	}

	return reply, nil

}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	movieapi.RegisterMovieInfoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
