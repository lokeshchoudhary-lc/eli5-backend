package question

import (
	"eli5/config/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func AskQuestion(c *fiber.Ctx) error {
	uniqueAlias := c.Locals("uniqueAlias").(string)
	askedQuestion := new(AskedQuestion)

	if err := c.BodyParser(askedQuestion); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	askedQuestion.AskedBy = uniqueAlias

	_, err := database.MG.Db.Collection("askedQuestions").InsertOne(c.Context(), askedQuestion)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(200)
}

func Explore(c *fiber.Ctx) error {
	//give list of all tags present
	query := bson.D{{Key: "choosen", Value: true}}
	cursor, err := database.MG.Db.Collection("tags").Find(c.Context(), query)
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
