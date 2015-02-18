package study

type Card struct {
	Title    string            `json:"title"`
	ImageUrl string            `json:"image_url"`
	Fields   map[string]string `json:"fields"`
	Notes    string            `json:"notes"`
}

func NewCard(title, imgUrl, notes string, fields map[string]string) *Card {
	return &Card{
		Title:    title,
		ImageUrl: imgUrl,
		Fields:   fields,
		Notes:    notes,
	}
}
