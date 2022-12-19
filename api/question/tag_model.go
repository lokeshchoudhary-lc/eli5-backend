package question

type Tag struct {
	Tag     string `json:"tag"`
	Choosen bool   `json:"choosen,omitempty" bson:"choosen,omitempty"`
}
