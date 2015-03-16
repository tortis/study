package main

import "gopkg.in/mgo.v2/bson"

type Deck struct {
	Id     bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string        `json:"name" bson:"name"`
	Fields []string      `json:"fields" bson:"fields"`
	Cards  []*Card       `json:"cards,omitempty" bson:"cards"`
}

func NewDeck(name string, fields []string) *Deck {
	return &Deck{Name: name, Fields: fields}
}
