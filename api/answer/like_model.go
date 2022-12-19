package answer

type Like struct {
	UserId   string `json:"userId" bson:"userId"`
	AnswerId string `json:"answerId" bson:"answerId"`
}
