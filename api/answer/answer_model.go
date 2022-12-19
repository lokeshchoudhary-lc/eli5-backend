package answer

type Answer struct {
	Id         string `json:"id,omitempty" bson:"_id,omitempty"`
	QuestionId string `json:"questionId" bson:"questionId"`
	Answer     string `json:"answer" bson:"answer"`
	AnsweredBy string `json:"answeredBy" bson:"answeredBy"`
	LikeNumber int64  `json:"likeNumber" bson:"likeNumber"`
	CreatedAt  int64  `json:"createdAt" bson:"createdAt"`
}
