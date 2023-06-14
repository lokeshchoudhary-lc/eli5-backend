package question

type Question struct {
	Id           string `json:"id,omitempty" bson:"_id,omitempty"`
	Question     string `json:"question" bson:"question"`
	Tag          string `json:"tag" bson:"tag"`
	Choosen      bool   `json:"choosen,omitempty" bson:"choosen,omitempty"`
	GifUrl       string `json:"gifUrl,omitempty" bson:"gifUrl,omitempty"`
	QuestionMark bool   `json:"questionMark,omitempty" bson:"questionMark,omitempty"`
}
