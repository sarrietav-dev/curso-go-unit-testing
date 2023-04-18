package main

import (
	"catching-pokemons/controller"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Add(a, b int) int {
	return a + b
}

func main() {
	router := chi.NewRouter()

	router.Get("/pokemon/{id}", controller.GetPokemon)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Print("Error found")
	}
}
