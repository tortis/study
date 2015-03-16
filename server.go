package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const PORT = ":8888"

var decks *mgo.Collection

func main() {
	// Create a gorilla router
	router := mux.NewRouter()

	// Open connection to database
	dbs, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	studyDB := dbs.DB("study")
	decks = studyDB.C("decks")

	// Attach handlers
	router.HandleFunc("/decks", listDecks).Methods("GET")
	router.HandleFunc("/decks", postDeck).Methods("POST")
	router.HandleFunc("/decks/{did}", getDeck).Methods("GET")
	router.HandleFunc("/decks/{did}", deleteDeck).Methods("DELETE")
	router.HandleFunc("/decks/{did}/cards", postCard).Methods("POST")
	router.HandleFunc("/decks/{did}/cards/{cid}", deleteCard).Methods("DELETE")
	router.HandleFunc("/decks/{did}/cards/{cid}", putCard).Methods("PUT")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	log.Printf("Starting server on port %s\n.", PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}

func listDecks(w http.ResponseWriter, r *http.Request) {
	// Response object that will be returned
	resp := struct {
		Decks []Deck
	}{}

	// Read decks from database, exclude the 'cards' field
	err := decks.Find(bson.M{}).Select(bson.M{"cards": 0}).All(&resp.Decks)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Marshel response to json and return
	j, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
func postDeck(w http.ResponseWriter, r *http.Request) {
	// Try to read a Deck from request body
	dec := json.NewDecoder(r.Body)
	d := Deck{}
	err := dec.Decode(&d)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Validate the deck
	if d.Name == "" {
		http.Error(w, "No name given for deck", http.StatusBadRequest)
		return
	}

	// Insert deck into the database
	err = decks.Insert(&d)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func getDeck(w http.ResponseWriter, r *http.Request) {
	deckId, err := parseObjectId(mux.Vars(r)["did"])
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	// Get the deck from the database
	d := Deck{}
	err = decks.FindId(deckId).One(&d)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
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
	// Get deck id from the url
	deckId, err := parseObjectId(mux.Vars(r)["did"])
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	// Delete deck from database
	err = decks.RemoveId(deckId)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func postCard(w http.ResponseWriter, r *http.Request) {
	// Get deck id from the url
	deckId, err := parseObjectId(mux.Vars(r)["did"])
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	// Read card from "json" part of body
	c := Card{}
	err = json.Unmarshal([]byte(r.FormValue("json")), &c)
	if err != nil {
		http.Error(w, "Failed to parse json.", http.StatusBadRequest)
		return
	}

	// Ensure a title was speicified.
	if c.Title == "" {
		http.Error(w, "The card must have a title.", http.StatusBadRequest)
		return
	}

	c.Id = bson.NewObjectId()

	// Read the file.
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Println("Could not open form file: ", err)
		http.Error(w, "There was a problem uploading the image.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate the file name
	fname := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)
	c.Image = fname

	// Write the file to disk
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

	// Add the new card to the deck
	q := bson.M{
		"$push": bson.M{"cards": &c},
	}
	err = decks.UpdateId(deckId, q)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteCard(w http.ResponseWriter, r *http.Request) {
	// Get deck id from the url
	deckId, err := parseObjectId(mux.Vars(r)["did"])
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	// Get the card number to delete
	cardId, err := parseObjectId(mux.Vars(r)["cid"])
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	q := bson.M{
		"$pull": bson.M{"cards": bson.M{"id": cardId}},
	}
	err = decks.UpdateId(deckId, q)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func putCard(w http.ResponseWriter, r *http.Request) {
	// Get deck id from the url
	deckId, err := parseObjectId(mux.Vars(r)["did"])
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	// Get the card number to delete
	cardId, err := parseObjectId(mux.Vars(r)["cid"])
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Read the new card from json body
	dec := json.NewDecoder(r.Body)
	c := Card{}
	err = dec.Decode(&c)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	// Validate the card
	if c.Title == "" {
		http.Error(w, "The card must have a title.", http.StatusBadRequest)
		return
	}

	// Update the card in the database
	q := bson.M{
		"$set": bson.M{
			"cards.$.title":  c.Title,
			"cards.$.fields": c.Fields,
			"cards.$.notes":  c.Notes,
		},
	}
	err = decks.Update(bson.M{"_id": deckId, "cards.id": cardId}, q)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func parseObjectId(s string) (bson.ObjectId, error) {
	if !bson.IsObjectIdHex(s) {
		return bson.NewObjectId(), fmt.Errorf("Id: %s is not a bson id.", s)
	}
	return bson.ObjectIdHex(s), nil
}
