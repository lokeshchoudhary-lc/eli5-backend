package feed

import (
	"eli5/api/question"
	"eli5/config/database"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MakeHomeFeed(c *fiber.Ctx) error {

	type answer struct {
		Id                 string `json:"id" bson:"_id"`
		QuestionId         string `json:"questionId" bson:"questionId"`
		Answer             string `json:"answer" bson:"answer"`
		AnsweredByName     string `json:"answeredByName" bson:"answeredByName"`
		AnsweredBy         string `json:"answeredBy" bson:"answeredBy"`
		LikeNumber         int64  `json:"likeNumber" bson:"likeNumber"`
		CreatedAt          int64  `json:"createdAt" bson:"createdAt"`
		Liked              bool   `json:"liked"`
		Tag                string `json:"tag" bson:"tag"`
		Question           string `json:"question" bson:"question"`
		ProfilePictureCode string `json:"profilePictureCode" bson:"profilePictureCode"`
	}

	var banners []Banner = make([]Banner, 0)
	var bestAnswer []answer = make([]answer, 0)
	var topTags []string
	var topQuestions []question.Question = make([]question.Question, 0)

	type makeTopQuestions struct {
		Id string `json:"id" bson:"_id"`
	}

	//find top questions Id's
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$questionId"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}}
	unsetStage := bson.D{{Key: "$unset", Value: bson.A{"count"}}}
	limitStage := bson.D{{Key: "$limit", Value: 5}}

	cursor, err := database.MG.Db.Collection("answers").Aggregate(c.Context(), mongo.Pipeline{groupStage, sortStage, unsetStage, limitStage})

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var results []makeTopQuestions = make([]makeTopQuestions, 0)

	if err := cursor.All(c.Context(), &results); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	//find questions for every Id form above
	//And take tag from them for tops tags
	for _, result := range results {
		var question question.Question
		id, _ := primitive.ObjectIDFromHex(result.Id)
		query := bson.D{{Key: "_id", Value: id}}
		err := database.MG.Db.Collection("questions").FindOne(c.Context(), query).Decode(&question)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(500).SendString(err.Error())
			}
			return c.Status(500).SendString(err.Error())
		}
		topQuestions = append(topQuestions, question)
		topTags = append(topTags, question.Tag)
	}

	//find the best answer with most like number
	opts := options.Find().SetSort(bson.D{{Key: "likeNumber", Value: -1}}).SetLimit(1)
	query := bson.D{{}}
	cursor, err = database.MG.Db.Collection("answers").Find(c.Context(), query, opts)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if err := cursor.All(c.Context(), &bestAnswer); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	bestAnswerSingle := bestAnswer[0]

	//find profile picture code of the user for the answer
	type user struct {
		ProfilePictureCode string `json:"profilePictureCode" bson:"profilePictureCode"`
		UniqueAlias        string `json:"uniqueAlias" bson:"uniqueAlias"`
	}
	var u user
	id, _ := primitive.ObjectIDFromHex(bestAnswerSingle.AnsweredBy)
	opt := options.FindOne().SetProjection(bson.D{{Key: "profilePictureCode", Value: 1}, {Key: "uniqueAlias", Value: 1}})
	query = bson.D{{Key: "_id", Value: id}}
	err = database.MG.Db.Collection("users").FindOne(c.Context(), query, opt).Decode(&u)
	if err != nil {

		if err == mongo.ErrNoDocuments {

			return c.Status(500).SendString(err.Error())
		}
		return c.Status(500).SendString(err.Error())
	}
	bestAnswerSingle.ProfilePictureCode = u.ProfilePictureCode
	bestAnswerSingle.AnsweredByName = u.UniqueAlias

	//find the question of the answer
	var question question.Question
	id, _ = primitive.ObjectIDFromHex(bestAnswerSingle.QuestionId)
	query = bson.D{{Key: "_id", Value: id}}
	err = database.MG.Db.Collection("questions").FindOne(c.Context(), query).Decode(&question)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return c.Status(500).SendString(err.Error())
		}
		return c.Status(500).SendString(err.Error())
	}
	bestAnswerSingle.Question = question.Question
	bestAnswerSingle.Tag = question.Tag

	//find the banners
	query = bson.D{{Key: "choosen", Value: true}}
	cursor, err = database.MG.Db.Collection("banners").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	if err := cursor.All(c.Context(), &banners); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(&fiber.Map{
		"banner":       banners,
		"topQuestions": topQuestions,
		"topTags":      topTags,
		"bestAnswer":   bestAnswerSingle,
	})
}
