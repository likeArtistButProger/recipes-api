package main

import (
	"log"
	"net/http"
	"test-api/api/recipe"
)

func main() {
	mainHandler := http.NewServeMux()

	mainHandler.Handle("/recipe/create", recipe.Handler)
	mainHandler.Handle("/recipe/get-list", recipe.Handler)
	mainHandler.Handle("/recipe/get-one-detailed", recipe.Handler)
	mainHandler.Handle("/recipe/update", recipe.Handler)
	mainHandler.Handle("/recipe/delete", recipe.Handler)

	err := http.ListenAndServe(":8000", mainHandler)
	if err != nil {
		log.Fatal(err)
	}
}
