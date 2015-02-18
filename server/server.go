package main

import (
	"fmt"
	"github.com/gorilla/pat"
	"github.com/tortis/study"
	"log"
	"net/http"
)

var dm *study.DeckManager

func main() {
	// Create a gorilla router
	router := pat.New()

	// Open the deck store
	var err error
	dm, err = study.Open("decks.gob")
	if err != nil {
		log.Fatal("Failed to open the decks gob file.")
	}

	// Attach handlers
	router.Get("/decks", listDecks)
	router.Post("/decks", postDeck)
	router.Get("/decks/{did}", getDeck)
	router.Delete("/decks/{did}", deleteDeck)
	router.Post("/decks/{did}/cards", postCard)
	router.Put("/decks/{did}/cards/{cid}", putCard)
	router.Delete("/decks/{did}/cards/{cid}", deleteCard)

	log.Printf("Starting server on port 8080.")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func listDecks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v", dm.ListDecks())
}
func postDeck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func getDeck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func deleteDeck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func postCard(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func deleteCard(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func putCard(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}
