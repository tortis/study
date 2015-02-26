package study

type Deck struct {
	Name   string   `json:"name"`
	Fields []string `json:"fields"`
	Cards  []*Card  `json:"cards"`
}

func NewDeck(name string, fields []string) *Deck {
	return &Deck{Name: name, Fields: fields}
}

func (d *Deck) AddField(name string) {
	d.Fields = append(d.Fields, name)
}

func (d *Deck) AddCard(c *Card) {
	d.Cards = append(d.Cards, c)
}
