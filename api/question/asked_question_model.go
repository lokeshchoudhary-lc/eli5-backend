package question

type AskedQuestion struct {
	Question string `json:"question" bson:"question"`
	AskedBy  string `json:"askedBy" bson:"askedBy"`
}
