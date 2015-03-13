package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

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
	router.HandleFunc("/decks/{did}/cards/{cn}", putCard).Methods("PUT")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	log.Printf("Starting server on port 8080.")
	log.Fatal(http.ListenAndServe(":8888", router))
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
	lock.Unlock()
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
	deckName, _ := url.QueryUnescape(mux.Vars(r)["did"])
	_, e := decks[deckName]
	if !e {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	delete(decks, deckName)
	err := saveDecks()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func postCard(w http.ResponseWriter, r *http.Request) {
	log.Println("Content-Type: ", r.Header.Get("Content-Type"))
	deckName, _ := url.QueryUnescape(mux.Vars(r)["did"])
	d, e := decks[deckName]
	if !e {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	c := study.Card{}
	err := json.Unmarshal([]byte(r.FormValue("json")), &c)
	if err != nil {
		http.Error(w, "Failed to parse json.", http.StatusBadRequest)
		return
	}

	// Ensure a title was speicified.
	if c.Title == "" {
		http.Error(w, "The card must have a title.", http.StatusBadRequest)
		return
	}

	// Read the file.
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Println("Could not open form file: ", err)
		http.Error(w, "There was a problem uploading the image.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fname := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)
	c.Image = fname
	out, err := os.OpenFile("static/images/"+fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Println("Failed to open image file", err)
		return
	}
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Println("Failed to copy file from http form.")
		return
	}

	d.Cards = append(d.Cards, &c)

	// Save state
	err = saveDecks()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteCard(w http.ResponseWriter, r *http.Request) {
	deckName, _ := url.QueryUnescape(mux.Vars(r)["did"])
	d, e := decks[deckName]
	if !e {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	cardNum, _ := strconv.Atoi(mux.Vars(r)["cn"])

	if cardNum >= len(d.Cards) || cardNum < 0 {
		http.Error(w, "Invalid card number", http.StatusBadRequest)
		return
	}

	// delete the card
	copy(d.Cards[cardNum:], d.Cards[cardNum+1:])
	d.Cards[len(d.Cards)-1] = nil
	d.Cards = d.Cards[:len(d.Cards)-1]

	// Save state
	err := saveDecks()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func putCard(w http.ResponseWriter, r *http.Request) {
	// Get the deck from url
	log.Println("Handing card update request.")
	deckName, _ := url.QueryUnescape(mux.Vars(r)["did"])
	d, e := decks[deckName]
	if !e {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	// Get the card number from url
	cardNum, _ := strconv.Atoi(mux.Vars(r)["cn"])
	if cardNum >= len(d.Cards) || cardNum < 0 {
		http.Error(w, "Invalid card number", http.StatusBadRequest)
		return
	}

	// Read the new card from json body
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

	d.Cards[cardNum] = &c
	// Save state
	err = saveDecks()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
