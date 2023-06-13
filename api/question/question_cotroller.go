package question

import (
	"eli5/config/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetQuestion(c *fiber.Ctx) error {
	questionId := c.Params("questionId")

	id, _ := primitive.ObjectIDFromHex(questionId)

	var question Question
	query := bson.D{{Key: "_id", Value: id}}
	err := database.MG.Db.Collection("questions").FindOne(c.Context(), query).Decode(&question)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(question)
}

func GetTrendingTags(c *fiber.Ctx) error {
	type topTags struct {
		Id string `json:"id" bson:"_id"`
	}

	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$tag"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}}
	unsetStage := bson.D{{Key: "$unset", Value: bson.A{"count"}}}
	limitStage := bson.D{{Key: "$limit", Value: 5}}

	cursor, err := database.MG.Db.Collection("questions").Aggregate(c.Context(), mongo.Pipeline{groupStage, sortStage, unsetStage, limitStage})

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var results []topTags = make([]topTags, 0)

	if err := cursor.All(c.Context(), &results); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(results)
}

func GetAllQuestions(c *fiber.Ctx) error {

	opts := options.Find()

	var perPageItem int64 = 3
	page, _ := strconv.Atoi(c.Query("page", "1"))
	opts.SetSort(bson.D{{Key: "_id", Value: -1}})
	opts.SetLimit(perPageItem)
	opts.SetSkip((int64(page) - 1) * perPageItem)

	query := bson.D{{}}

	cursor, err := database.MG.Db.Collection("questions").Find(c.Context(), query, opts)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	var questions []Question = make([]Question, 0)

	if err := cursor.All(c.Context(), &questions); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if length := len(questions); length == 0 {
		return c.SendStatus(204)
	}
	return c.JSON(questions)
}

func AskQuestion(c *fiber.Ctx) error {
	askedQuestion := new(AskedQuestion)

	if err := c.BodyParser(askedQuestion); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	_, err := database.MG.Db.Collection("askedQuestions").InsertOne(c.Context(), askedQuestion)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(200)
}

func TagsPageStats(c *fiber.Ctx) error {
	tag := c.Params("tag")

	type count struct {
		LikeCount   int `json:"likeCount" bson:"likeCount"`
		AnswerCount int `json:"answerCount" bson:"answerCount"`
	}

	var tagCount []count

	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "tag", Value: tag},
		}}}
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$tag"},
			{Key: "likeCount", Value: bson.D{{Key: "$sum", Value: "$likeNumber"}}},
			{Key: "answerCount", Value: bson.D{{Key: "$count", Value: bson.D{}}}},
		}}}

	cursor, err := database.MG.Db.Collection("answers").Aggregate(c.Context(), mongo.Pipeline{matchStage, groupStage})

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if err := cursor.All(c.Context(), &tagCount); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if len(tagCount) > 0 {
		return c.JSON(tagCount[0])
	}
	return c.SendStatus(200)
}

func Explore(c *fiber.Ctx) error {
	//give list of all tags present
	options := options.Find()
	options.SetSort(bson.D{{Key: "_id", Value: -1}})
	query := bson.D{{Key: "choosen", Value: true}}
	cursor, err := database.MG.Db.Collection("tags").Find(c.Context(), query, options)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	var allTags []Tag = make([]Tag, 0)

	if err := cursor.All(c.Context(), &allTags); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(allTags)
}

func GetTags(c *fiber.Ctx) error {
	opts := options.Find()

	var perPageItem int64 = 10
	page, _ := strconv.Atoi(c.Query("page", "1"))
	opts.SetLimit(perPageItem)
	opts.SetSkip((int64(page) - 1) * perPageItem)

	query := bson.D{{Key: "choosen", Value: true}}
	cursor, err := database.MG.Db.Collection("tags").Find(c.Context(), query, opts)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	var allTags []Tag = make([]Tag, 0)

	if err := cursor.All(c.Context(), &allTags); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	if length := len(allTags); length == 0 {
		return c.SendStatus(204)
	}
	return c.JSON(allTags)
}

func GetQuestionsOfTag(c *fiber.Ctx) error {
	tag := c.Params("tag")

	opts := options.Find()

	var perPageItem int64 = 10
	page, _ := strconv.Atoi(c.Query("page", "1"))
	opts.SetLimit(perPageItem)
	opts.SetSkip((int64(page) - 1) * perPageItem)

	query := bson.D{{Key: "$and", Value: bson.A{bson.D{{Key: "choosen", Value: true}}, bson.D{{Key: "tag", Value: tag}}}}}

	cursor, err := database.MG.Db.Collection("questions").Find(c.Context(), query, opts)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	var questions []Question = make([]Question, 0)

	if err := cursor.All(c.Context(), &questions); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if length := len(questions); length == 0 {
		return c.SendStatus(204)
	}
	return c.JSON(questions)
}

func GetBacklinkQuestionsOfTag(c *fiber.Ctx) error {
	tag := c.Params("tag")

	opts := options.Find()

	var perPageItem int64 = 10
	opts.SetSort(bson.D{{Key: "_id", Value: -1}})
	opts.SetLimit(perPageItem)

	query := bson.D{{Key: "$and", Value: bson.A{bson.D{{Key: "choosen", Value: true}}, bson.D{{Key: "tag", Value: tag}}}}}

	cursor, err := database.MG.Db.Collection("questions").Find(c.Context(), query, opts)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	var questions []Question = make([]Question, 0)

	if err := cursor.All(c.Context(), &questions); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(questions)
}

// func GetQuestionOfTag(c *fiber.Ctx) error {
// 	tag := c.Params("tag")

// 	var question Question
// 	query := bson.D{{Key: "$and", Value: bson.A{bson.D{{Key: "choosen", Value: true}}, bson.D{{Key: "tag", Value: tag}}}}}
// 	err := database.MG.Db.Collection("questions").FindOne(c.Context(), query).Decode(&question)

// 	if err != nil {

// 		return c.Status(500).SendString(err.Error())
// 	}
// 	return c.JSON(question)

// }

func ChangeQuestionWithPagination(c *fiber.Ctx) error {
	questionId := c.Params("questionId")
	tag := c.Query("tag")

	id, _ := primitive.ObjectIDFromHex(questionId)

	opts := options.Find()

	var perPageItem int64 = 10
	page, _ := strconv.Atoi(c.Query("page", "1"))
	opts.SetLimit(perPageItem)
	opts.SetSkip((int64(page) - 1) * perPageItem)

	if c.Query("action") == "next" {

		query := bson.D{{Key: "$and", Value: bson.A{
			bson.D{{Key: "_id", Value: bson.D{{Key: "$gt", Value: id}}}},
			bson.D{{Key: "$and", Value: bson.A{
				bson.D{{Key: "choosen", Value: true}}, bson.D{{Key: "tag", Value: tag}}}}}}}}

		cursor, err := database.MG.Db.Collection("questions").Find(c.Context(), query, opts)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		var questions []Question = make([]Question, 0)

		if err := cursor.All(c.Context(), &questions); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		if length := len(questions); length == 0 {
			return c.SendStatus(204)
		}

		return c.JSON(questions)

	}

	if c.Query("action") == "back" {

		opts.SetSort(bson.D{{Key: "_id", Value: -1}})

		query := bson.D{{Key: "$and", Value: bson.A{bson.D{{Key: "_id", Value: bson.D{{Key: "$lt", Value: id}}}}, bson.D{{Key: "$and", Value: bson.A{bson.D{{Key: "choosen", Value: true}}, bson.D{{Key: "tag", Value: tag}}}}}}}}
		cursor, err := database.MG.Db.Collection("questions").Find(c.Context(), query, opts)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		var questions []Question = make([]Question, 0)

		if err := cursor.All(c.Context(), &questions); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		if length := len(questions); length == 0 {
			return c.SendStatus(204)
		}
		return c.JSON(questions)

	}
	return c.SendStatus(200)
}
