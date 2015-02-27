package study

type Card struct {
	Title  string            `json:"title"`
	Image  string            `json:"image"`
	Fields map[string]string `json:"fields"`
	Notes  string            `json:"notes"`
}

func NewCard(title, imgName, notes string, fields map[string]string) *Card {
	return &Card{
		Title:  title,
		Image:  imgName,
		Fields: fields,
		Notes:  notes,
	}
}
