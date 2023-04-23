package company

type Company struct {
	Id                string `json:"id,omitempty" bson:"_id,omitempty"`
	UniqueAlias       string `json:"uniqueAlias,omitempty" bson:"uniqueAlias,omitempty"`
	Name              string `json:"name,omitempty" bson:"name,omitempty"`
	Bio               string `json:"bio,omitempty" bson:"bio,omitempty"`
	Email             string `json:"email,omitempty" bson:"email,omitempty"`
	BrandColor        string `json:"brandColor,omitempty" bson:"brandColor,omitempty"`
	WebsiteUrl        string `json:"websiteUrl,omitempty" bson:"websiteUrl,omitempty"`
	TwitterUrl        string `json:"twitterUrl,omitempty" bson:"twitterUrl,omitempty"`
	LinkedinUrl       string `json:"linkedinUrl,omitempty" bson:"linkedinUrl,omitempty"`
	FacebookUrl       string `json:"facebookUrl,omitempty" bson:"facebookUrl,omitempty"`
	YoutubeUrl        string `json:"youtubeUrl,omitempty" bson:"youtubeUrl,omitempty"`
	InstagramUrl      string `json:"instagramUrl,omitempty" bson:"instagramUrl,omitempty"`
	OtherUrl          string `json:"otherUrl,omitempty" bson:"otherUrl,omitempty"`
	ProfilePictureUrl string `json:"profilePictureUrl,omitempty" bson:"profilePictureUrl,omitempty"`
	TotalLikes        int64  `json:"totalLikes" bson:"totalLikes"`
	CreatedAt         int64  `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
