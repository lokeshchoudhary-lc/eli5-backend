package user

type User struct {
	Id                 string `json:"id,omitempty" bson:"_id,omitempty"`
	UniqueAlias        string `json:"uniqueAlias,omitempty" bson:"uniqueAlias,omitempty"`
	FirstName          string `json:"firstName,omitempty" bson:"firstName,omitempty"`
	Bio                string `json:"bio,omitempty" bson:"bio,omitempty"`
	Email              string `json:"email,omitempty" bson:"email,omitempty"`
	TwitterUrl         string `json:"twitterUrl,omitempty" bson:"twitterUrl,omitempty"`
	LinkedinUrl        string `json:"linkedinUrl,omitempty" bson:"linkedinUrl,omitempty"`
	InstagramUrl       string `json:"instagramUrl,omitempty" bson:"instagramUrl,omitempty"`
	ReferedBy          string `json:"referedBy,omitempty" bson:"referedBy,omitempty"`
	ProfilePictureCode string `json:"profilePictureCode,omitempty" bson:"profilePictureCode,omitempty"`
	Streak             int64  `json:"streak" bson:"streak"`
	TotalAnswers       int64  `json:"totalAnswers" bson:"totalAnswers"`
	TotalLikes         int64  `json:"totalLikes" bson:"totalLikes"`
	CreatedAt          int64  `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	PrevStreakTime     int64  `json:"prevStreakTime,omitempty" bson:"prevStreakTime,omitempty"`
	// Score              int64  `json:"score" bson:"score"`
}

//make index on uniqueAlias (future)
// make index on total likes and answer , descending
