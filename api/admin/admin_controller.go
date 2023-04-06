package admin

import (
	"eli5/api/answer"
	"eli5/api/question"
	"eli5/config/database"
	"eli5/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTags(c *fiber.Ctx) error {
	query := bson.D{{}}
	cursor, err := database.MG.Db.Collection("tags").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	var tags []string
	if err := cursor.All(c.Context(), &tags); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(tags)

}
func GetQuestionsOfTag(c *fiber.Ctx) error {
	tagName := c.Params("tagName")
	query := bson.D{{Key: "tag", Value: tagName}}
	cursor, err := database.MG.Db.Collection("questions").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	var questions []string
	if err := cursor.All(c.Context(), &questions); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(questions)
}
func ChooseQuestion(c *fiber.Ctx) error {
	id, _ := primitive.ObjectIDFromHex(c.Params("questionId"))
	query := bson.D{{Key: "_id", Value: id}}
	updateQuery := bson.D{{Key: "$set", Value: bson.D{{Key: "choosen", Value: true}}}}
	_, err := database.MG.Db.Collection("questions").UpdateOne(c.Context(), query, updateQuery)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Status(200).SendString("Successfully Choosen")
}
func ChooseTag(c *fiber.Ctx) error {
	tagName := c.Params("tagName")
	query := bson.D{{Key: "tag", Value: tagName}}
	updateQuery := bson.D{{Key: "$set", Value: bson.D{{Key: "choosen", Value: true}}}}
	_, err := database.MG.Db.Collection("tags").UpdateOne(c.Context(), query, updateQuery)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Status(200).SendString("Successfully Choosen")
}
func PostQuestion(c *fiber.Ctx) error {
	newQuestion := new(question.Question)

	if err := c.BodyParser(newQuestion); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	newQuestion.Id = ""

	insertResult, err := database.MG.Db.Collection("questions").InsertOne(c.Context(), newQuestion)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	query := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	createdRecord := database.MG.Db.Collection("questions").FindOne(c.Context(), query)

	createdQuestion := &question.Question{}
	createdRecord.Decode(createdQuestion)

	return c.Status(201).JSON(createdQuestion)

}
func PostTag(c *fiber.Ctx) error {
	tagName := c.Params("tagName")
	newTag := new(question.Tag)
	newTag.Tag = tagName

	insertResult, err := database.MG.Db.Collection("tags").InsertOne(c.Context(), newTag)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	query := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	createdRecord := database.MG.Db.Collection("tags").FindOne(c.Context(), query)

	createdTag := &question.Tag{}
	createdRecord.Decode(createdTag)

	return c.Status(201).JSON(createdTag)
}

func PostGptAnswer(c *fiber.Ctx) error {

	gptAnswer := new(answer.GptAnswer)

	if err := c.BodyParser(gptAnswer); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// force MongoDB to always set its own generated ObjectIDs
	gptAnswer.Id = ""

	//sanitize answer
	html := utils.SanitizeHtml(gptAnswer.Answer)
	gptAnswer.Answer = html

	// insert the record
	insertResult, err := database.MG.Db.Collection("gptAnswers").InsertOne(c.Context(), gptAnswer)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	createdGptAnswer := new(answer.GptAnswer)

	// get the just inserted record in order to return it as response
	query := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	database.MG.Db.Collection("gptAnswers").FindOne(c.Context(), query).Decode(&createdGptAnswer)

	return c.Status(201).JSON(&fiber.Map{"gptAnswer": createdGptAnswer})

}

// func UpdateQuestion(c *fiber.Ctx) error {
// 	id, _ := primitive.ObjectIDFromHex(c.Params("questionId"))
// 	query := bson.D{{Key: "_id", Value: id}}
// 	updateQuery := bson.D{{Key: "$set", Value: bson.D{{Key: "choosen", Value: true}}}}
// 	_, err := database.MG.Db.Collection("questions").UpdateOne(c.Context(), query, updateQuery)
// 	if err != nil {
// 		return c.Status(500).SendString(err.Error())
// 	}
// 	return c.Status(200).SendString("Successfully Choosen")
// }

// func DeleteQuestion(c *fiber.Ctx) error {

// }
// func DeleteTag(c *fiber.Ctx) error {

// }
