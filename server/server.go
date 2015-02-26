package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"github.com/tortis/study"
)

var decks map[string]*study.Deck
var lock sync.Mutex

// Load decks from json if possible
func init() {
	lock = sync.Mutex{}
	f, err := os.Open("decks.json")
	if err != nil {
		log.Println("No deck.json file found. Starting fresh")
		decks = make(map[string]*study.Deck)
		return
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(&decks)
	if err != nil {
		log.Println("Failed to read deck.json file. Starting fresh")
		decks = make(map[string]*study.Deck)
	}
}

func main() {
	// Create a gorilla router
	router := mux.NewRouter()

	// Attach handlers
	router.HandleFunc("/decks", listDecks).Methods("GET")
	router.HandleFunc("/decks", postDeck).Methods("POST")
	router.HandleFunc("/decks/{did}", getDeck).Methods("GET")
	router.HandleFunc("/decks/{did}", deleteDeck).Methods("DELETE")
	router.HandleFunc("/decks/{did}/cards", postCard).Methods("POST")
	router.HandleFunc("/decks/{did}/cards/{cn}", deleteCard).Methods("DELETE")
	router.HandleFunc("/decks/{did}/cards/{cn}", putCard).Methods("DELETE")

	log.Printf("Starting server on port 8080.")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func saveDecks() error {
	lock.Lock()
	f, err := os.Create("decks.json")
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	err = enc.Encode(&decks)
	if err != nil {
		return err
	}
	return nil
}

func listDecks(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Decks []string
	}{Decks: make([]string, 0)}
	for n, _ := range decks {
		resp.Decks = append(resp.Decks, n)
	}
	j, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
func postDeck(w http.ResponseWriter, r *http.Request) {
	log.Println("postDeck handler called.")
	// Try to read a Deck from request body
	dec := json.NewDecoder(r.Body)
	d := study.Deck{}
	err := dec.Decode(&d)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	decks[d.Name] = &d
	err = saveDecks()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func getDeck(w http.ResponseWriter, r *http.Request) {
	urlName := mux.Vars(r)["did"]
	name, _ := url.QueryUnescape(urlName)
	log.Printf("Searching for %s\n", name)
	d, e := decks[name]
	if !e {
		http.Error(w, "Deck not found", http.StatusNotFound)
		return
	}
	b, err := json.Marshal(&d)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
func deleteDeck(w http.ResponseWriter, r *http.Request) {
	log.Printf("Deleting deck: %s\n", mux.Vars(r)["did"])
	http.Error(w, "", http.StatusNotImplemented)
}
func postCard(w http.ResponseWriter, r *http.Request) {
	log.Println("postCard handler called")
	// Ensure the deck exists.
	deckName, _ := url.QueryUnescape(mux.Vars(r)["did"])
	_, e := decks[deckName]
	if !e {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	dec := json.NewDecoder(r.Body)
	c := study.Card{}
	err := dec.Decode(&c)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	// Validate the card
	if c.Title == "" {
		http.Error(w, "The card must have a title.", http.StatusBadRequest)
		return
	}
	// Add the card to the deck
	decks[deckName].Cards = append(decks[deckName].Cards, &c)

	// Save state
	err = saveDecks()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func deleteCard(w http.ResponseWriter, r *http.Request) {
	deckName, _ := url.QueryUnescaped(mux.Vars(r)["did"])
	_, e := decks[deckName]
	if !e {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	cardNum, _ := strconv.Atoi(mux.Vars(r)["cn"])
	decks[deckName].Cards = decks[deckName].Cards[
	w.WriteHeader(http.StatusOK)
}
func putCard(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
