package seed

import (
	"context"
	"eli5/api/question"
	"eli5/config/database"
)

func SeedDatabase() {
	tag := new(question.Tag)
	tag1 := new(question.Tag)
	tag2 := new(question.Tag)
	tag3 := new(question.Tag)
	tag4 := new(question.Tag)
	q1 := new(question.Question)
	q2 := new(question.Question)
	q3 := new(question.Question)
	q4 := new(question.Question)
	q5 := new(question.Question)

	tag.Tag = "computer"
	tag.Choosen = true
	tag1.Tag = "golang"
	tag1.Choosen = true
	tag2.Tag = "javascript"
	tag2.Choosen = true
	tag3.Tag = "cpp"
	tag3.Choosen = true
	tag4.Tag = "php"
	tag4.Choosen = true

	q1.Id = ""
	q1.Question = "What is a computer"
	q1.Choosen = true
	q1.Tag = "computer"

	q2.Question = "What is a golang"
	q2.Choosen = true
	q2.Tag = "golang"
	q3.Id = ""
	q3.Question = "What is a javascript"
	q3.Choosen = true
	q3.Tag = "javascript"
	q4.Id = ""
	q4.Question = "What is a cpp"
	q4.Choosen = true
	q4.Tag = "cpp"
	q5.Id = ""
	q5.Question = "What is a php"
	q5.Choosen = true
	q5.Tag = "php"

	// database.MG.Db.Collection("questions").Drop(context.Background())
	// database.MG.Db.Collection("tags").Drop(context.Background())

	database.MG.Db.Collection("tags").InsertOne(context.TODO(), tag)
	database.MG.Db.Collection("tags").InsertOne(context.TODO(), tag1)
	database.MG.Db.Collection("tags").InsertOne(context.TODO(), tag2)
	database.MG.Db.Collection("tags").InsertOne(context.TODO(), tag3)
	database.MG.Db.Collection("tags").InsertOne(context.TODO(), tag4)

	database.MG.Db.Collection("questions").InsertOne(context.TODO(), q1)
	database.MG.Db.Collection("questions").InsertOne(context.TODO(), q2)
	database.MG.Db.Collection("questions").InsertOne(context.TODO(), q3)
	database.MG.Db.Collection("questions").InsertOne(context.TODO(), q4)
	database.MG.Db.Collection("questions").InsertOne(context.TODO(), q5)

}
