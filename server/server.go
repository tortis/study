package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/pat"
	"github.com/tortis/study"
)

var dm *study.DeckManager

func main() {
	// Create a gorilla router
	router := pat.New()

	// Open the deck store
	var err error
	dm, err = study.Open("decks.gob")
	if err != nil {
		log.Fatal("Failed to open the decks gob .")
	}

	// Add a test deck
	dm.AddDeck(study.NewDeck("Test Deck"))

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
	resp := struct {
		Decks []string
	}{Decks: dm.ListDecks()}
	j, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Write(j)
}
func postDeck(w http.ResponseWriter, r *http.Request) {
	var d study.Deck
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&d)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	dm.AddDeck(&d)
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
