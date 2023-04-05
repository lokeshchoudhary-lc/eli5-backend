package answer

type GptAnswer struct {
	Id         string `json:"id,omitempty" bson:"_id,omitempty"`
	QuestionId string `json:"questionId" bson:"questionId"`
	Answer     string `json:"answer" bson:"answer"`
	// AnsweredById string `json:"answeredById" bson:"answeredById"`
}
