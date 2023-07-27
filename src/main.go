package main

import (
	"log"
	"net/http"
	"test-api/api/recipe"
)

func main() {
	mainHandler := http.NewServeMux()
	mainHandler.Handle("/recipe", recipe.Handler)

	err := http.ListenAndServe(":8000", mainHandler)
	if err != nil {
		log.Fatal(err)
	}
}
