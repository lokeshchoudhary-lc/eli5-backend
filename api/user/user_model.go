package user

type User struct {
	Id                 string `json:"id,omitempty" bson:"_id,omitempty"`
	UniqueAlias        string `json:"uniqueAlias,omitempty" bson:"uniqueAlias,omitempty"`
	FirstName          string `json:"firstName,omitempty" bson:"firstName,omitempty"`
	Email              string `json:"email,omitempty" bson:"email,omitempty"`
	ProfilePictureCode string `json:"profilePictureCode,omitempty" bson:"profilePictureCode,omitempty"`
	Streak             int64  `json:"streak" bson:"streak"`
	// Score              int64  `json:"score" bson:"score"`
	TotalAnswers   int64 `json:"totalAnswers" bson:"totalAnswers"`
	TotalLikes     int64 `json:"totalLikes" bson:"totalLikes"`
	CreatedAt      int64 `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	PrevStreakTime int64 `json:"prevStreakTime,omitempty" bson:"prevStreakTime,omitempty"`
}

//make index on uniqueAlias (future)
// make index on total likes and answer , descending
