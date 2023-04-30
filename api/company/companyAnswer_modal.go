package company

type CompanyAnswer struct {
	Id         string `json:"id,omitempty" bson:"_id,omitempty"`
	QuestionId string `json:"questionId" bson:"questionId"`
	Answer     string `json:"answer" bson:"answer"`
	CreatedAt  int64  `json:"createdAt" bson:"createdAt"`
}
