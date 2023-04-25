package company

type CompanyQuestion struct {
	Id        string `json:"id,omitempty" bson:"_id,omitempty"`
	Question  string `json:"question" bson:"question"`
	CompanyId string `json:"companyId" bson:"companyId"`
	Tag       string `json:"tag" bson:"tag"`
	Choosen   bool   `json:"choosen,omitempty" bson:"choosen,omitempty"`
}
