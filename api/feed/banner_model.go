package feed

type Banner struct {
	ImageUrl    string `json:"imageUrl" bson:"imageUrl"`
	ContentText string `json:"contentText" bson:"contentText"`
	Choosen     bool   `json:"choosen" bson:"choosen"`
	CTA         string `json:"cta" bson:"cta"`
}
