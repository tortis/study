package study

import (
	"encoding/gob"
	"io"
	"os"
)

type DeckManager struct {
	deckFile string
	decks    map[string]*Deck
}

func Open(deckFile string) (*DeckManager, error) {
	dm := &DeckManager{
		deckFile: deckFile,
		decks:    make(map[string]*Deck),
	}
	f, err := os.OpenFile(deckFile, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	for err == nil {
		var d Deck
		err = dec.Decode(&d)
		if err == nil {
			dm.decks[d.Name] = &d
		}
	}
	if err != io.EOF {
		return nil, err
	}
	return dm, nil
}

func (dm *DeckManager) AddDeck(d *Deck) {
	dm.decks[d.Name] = d
}

func (dm *DeckManager) GetDeck(name string) *Deck {
	return dm.decks[name]
}

func (dm *DeckManager) ListDecks() []string {
	l := make([]string, len(dm.decks))
	for name, _ := range dm.decks {
		l = append(l, name)
	}
	return l
}

func (dm *DeckManager) Save() error {
	f, err := os.Create(dm.deckFile)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	for _, deck := range dm.decks {
		err = enc.Encode(deck)
		if err != nil {
			return err
		}
	}
	return nil
}
