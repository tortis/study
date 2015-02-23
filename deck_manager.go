package study

import (
	"encoding/json"

	"github.com/peterbourgon/diskv"
)

type DeckManager struct {
	d *diskv.Diskv
}

func Open(dataPath string) (*DeckManager, error) {
	dm := &DeckManager{
		d: diskv.New(diskv.Options{
			BasePath:     dataPath,
			Transform:    func(s string) []string { return []string{} },
			CacheSizeMax: 1024 * 1024,
		}),
	}
	return dm, nil
}

func (dm *DeckManager) AddDeck(d *Deck) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	err = dm.d.Write(d.Name, b)
	if err != nil {
		return err
	}
	return nil
}

func (dm *DeckManager) RemoveDeck(d *Deck) {
}

func (dm *DeckManager) GetDeck(name string) (*Deck, error) {
	b, err := dm.d.Read(name)
	if err != nil {
		return nil, err
	}
	var d Deck
	err = json.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}
	return &d, err
}

func (dm *DeckManager) ListDecks() []string {
	l := make([]string, 0)
	k := dm.d.Keys(nil)
	for name := range k {
		l = append(l, name)
	}
	return l
}
