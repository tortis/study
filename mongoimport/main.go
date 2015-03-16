package main

import (
	"encoding/json"
	"log"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Card struct {
	Id     bson.ObjectId     `json:"id", bson:"id"`
	Title  string            `json:"title" bson:"title"`
	Image  string            `json:"image" bson:"image"`
	Fields map[string]string `json:"fields" bson:"fields"`
	Notes  string            `json:"notes" bson:"notes"`
}

type Deck struct {
	Id     bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string        `json:"name" bson:"name"`
	Fields []string      `json:"fields" bson:"fields"`
	Cards  []*Card       `json:"cards,omitempty" bson:"cards"`
}

func main() {
	// Load decks from json file
	f, err := os.Open("decks.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var decks map[string]*Deck
	err = dec.Decode(&decks)
	if err != nil {
		log.Fatal(err)
	}

	// Open mongo connection
	dbs, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	deckCol := dbs.DB("study").C("decks")

	// Import each deck, adding ids to deck and cards
	for _, deck := range decks {
		log.Printf("Import deck: %s\n", deck.Name)
		// Generate an ID for each card
		for _, c := range deck.Cards {
			c.Id = bson.NewObjectId()
		}
		err = deckCol.Insert(deck)
		if err != nil {
			log.Println(err)
		}
	}
}
