package main

import "gopkg.in/mgo.v2/bson"

type Card struct {
	Id     bson.ObjectId     `json:"id", bson:"id"`
	Title  string            `json:"title" bson:"title"`
	Image  string            `json:"image" bson:"image"`
	Fields map[string]string `json:"fields" bson:"fields"`
	Notes  string            `json:"notes" bson:"notes"`
}

func NewCard(title, imgName, notes string, fields map[string]string) *Card {
	return &Card{
		Id:     bson.NewObjectId(),
		Title:  title,
		Image:  imgName,
		Fields: fields,
		Notes:  notes,
	}
}
