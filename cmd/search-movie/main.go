package main

import (
	"fmt"
	"net/http"

	"github.com/kylecain/wheel-of-wonder/internal/service"
)

func main() {
	service := service.NewMovieSearch(http.DefaultClient)

	movie, err := service.FetchMovie("Inception")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", movie.Title)
	fmt.Printf("%s", movie.Description)
	fmt.Printf("%s", movie.ImageURL)
}
