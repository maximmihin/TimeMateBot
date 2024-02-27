package entities

type Tag struct {
	ID      int    `json:"id" csv:"ID"`
	Tag     string `json:"tag" csv:"Tag"`
	Comment string `json:"comment,omitempty" csv:"Comment,omitempty"`
}
