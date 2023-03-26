package answer

import (
	"eli5/config/database"
	"eli5/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collection string = "answers"

func GetUserAnswer(c *fiber.Ctx) error {
	var uniqueAlias string = c.Locals("uniqueAlias").(string)

	type answer struct {
		Id         string `json:"id" bson:"_id"`
		QuestionId string `json:"questionId" bson:"questionId"`
		Answer     string `json:"answer" bson:"answer"`
		AnsweredBy string `json:"answeredBy" bson:"answeredBy"`
		LikeNumber int64  `json:"likeNumber" bson:"likeNumber"`
		CreatedAt  int64  `json:"createdAt" bson:"createdAt"`
		Liked      bool   `json:"liked"`
	}
	var userAnswer answer

	query := bson.D{{Key: "$and", Value: bson.A{bson.D{{Key: "answeredBy", Value: uniqueAlias}}, bson.D{{Key: "questionId", Value: c.Params("questionId")}}}}}
	err := database.MG.Db.Collection(collection).FindOne(c.Context(), query).Decode(&userAnswer)

	if err == mongo.ErrNoDocuments {
		return c.JSON(&fiber.Map{"message": "no_answer"})

	}
	return c.JSON(userAnswer)

}

func GetGuestAnswers(c *fiber.Ctx) error {

	var query bson.D
	opts := options.Find()

	// questionId, _ := primitive.ObjectIDFromHex(c.Params("questionId"))
	questionId := c.Params("questionId")

	if sortOrder := c.Query("sort"); sortOrder == "lastest" {
		query = bson.D{{Key: "questionId", Value: questionId}}
		opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})
	}
	if sortOrder := c.Query("sort"); sortOrder == "trending" {
		query = bson.D{{Key: "questionId", Value: questionId}}
		opts.SetSort(bson.D{{Key: "likeNumber", Value: -1}})
	}
	//pagination
	var perPageItem int64 = 10
	page, _ := strconv.Atoi(c.Query("page", "1"))
	opts.SetLimit(perPageItem)
	opts.SetSkip((int64(page) - 1) * perPageItem)

	cursor, err := database.MG.Db.Collection(collection).Find(c.Context(), query, opts)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	type answer struct {
		Id                 string `json:"id" bson:"_id"`
		QuestionId         string `json:"questionId" bson:"questionId"`
		Answer             string `json:"answer" bson:"answer"`
		AnsweredBy         string `json:"answeredBy" bson:"answeredBy"`
		LikeNumber         int64  `json:"likeNumber" bson:"likeNumber"`
		CreatedAt          int64  `json:"createdAt" bson:"createdAt"`
		Liked              bool   `json:"liked"`
		ProfilePictureCode string `json:"profilePictureCode" bson:"profilePictureCode"`
	}

	var answers []answer = make([]answer, 0)

	if err := cursor.All(c.Context(), &answers); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if length := len(answers); length == 0 {
		return c.SendStatus(204)
	}

	// toggle like for each answer (only toggle like when user logged in) and find user profilePictureCode
	for i, answer := range answers {
		//give every answer liked as false for guest render of answers
		answers[i].Liked = false

		type ppc struct {
			ProfilePictureCode string `json:"profilePictureCode" bson:"profilePictureCode"`
		}
		var pp ppc
		opts := options.FindOne().SetProjection(bson.D{{Key: "profilePictureCode", Value: 1}})
		query = bson.D{{Key: "uniqueAlias", Value: answer.AnsweredBy}}
		err := database.MG.Db.Collection("users").FindOne(c.Context(), query, opts).Decode(&pp)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		answers[i].ProfilePictureCode = pp.ProfilePictureCode
	}
	return c.JSON(answers)

}
func GetAnswers(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)

	var query bson.D
	opts := options.Find()

	// questionId, _ := primitive.ObjectIDFromHex(c.Params("questionId"))
	questionId := c.Params("questionId")

	if sortOrder := c.Query("sort"); sortOrder == "lastest" {
		query = bson.D{{Key: "questionId", Value: questionId}}
		opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})
	}
	if sortOrder := c.Query("sort"); sortOrder == "trending" {
		query = bson.D{{Key: "questionId", Value: questionId}}
		opts.SetSort(bson.D{{Key: "likeNumber", Value: -1}})
	}
	//pagination
	var perPageItem int64 = 10
	page, _ := strconv.Atoi(c.Query("page", "1"))
	opts.SetLimit(perPageItem)
	opts.SetSkip((int64(page) - 1) * perPageItem)

	cursor, err := database.MG.Db.Collection(collection).Find(c.Context(), query, opts)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	type answer struct {
		Id                 string `json:"id" bson:"_id"`
		QuestionId         string `json:"questionId" bson:"questionId"`
		Answer             string `json:"answer" bson:"answer"`
		AnsweredBy         string `json:"answeredBy" bson:"answeredBy"`
		LikeNumber         int64  `json:"likeNumber" bson:"likeNumber"`
		CreatedAt          int64  `json:"createdAt" bson:"createdAt"`
		Liked              bool   `json:"liked"`
		ProfilePictureCode string `json:"profilePictureCode" bson:"profilePictureCode"`
	}

	var answers []answer = make([]answer, 0)

	if err := cursor.All(c.Context(), &answers); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if length := len(answers); length == 0 {
		return c.SendStatus(204)
	}

	// toggle like for each answer (only toggle like when user logged in) and find user profilePictureCode
	for i, answer := range answers {
		//should i use concurrency here ?

		query := bson.D{{Key: "$and", Value: bson.A{bson.D{{Key: "userId", Value: userId}}, bson.D{{Key: "answerId", Value: answer.Id}}}}}
		result := database.MG.Db.Collection("likes").FindOne(c.Context(), query)
		if result.Err() == nil {
			answers[i].Liked = true
		} else {
			answers[i].Liked = false
		}

		type ppc struct {
			ProfilePictureCode string `json:"profilePictureCode" bson:"profilePictureCode"`
		}
		var pp ppc
		opts := options.FindOne().SetProjection(bson.D{{Key: "profilePictureCode", Value: 1}})
		query = bson.D{{Key: "uniqueAlias", Value: answer.AnsweredBy}}
		err := database.MG.Db.Collection("users").FindOne(c.Context(), query, opts).Decode(&pp)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		answers[i].ProfilePictureCode = pp.ProfilePictureCode
	}
	return c.JSON(answers)

}

func PostAnswer(c *fiber.Ctx) error {

	//userId and uniqueAlias from token
	var userId string = c.Locals("userId").(string)
	var uniqueAlias string = c.Locals("uniqueAlias").(string)

	//check if already answered by the user uniqueAlias

	query := bson.D{{Key: "$and", Value: bson.A{bson.D{{Key: "answeredBy", Value: uniqueAlias}}, bson.D{{Key: "questionId", Value: c.Params("questionId")}}}}}

	res := database.MG.Db.Collection("answers").FindOne(c.Context(), query)

	if res.Err() != mongo.ErrNoDocuments {
		return c.SendStatus(400)
	}

	id, _ := primitive.ObjectIDFromHex(userId)
	answer := new(Answer)

	if err := c.BodyParser(answer); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	answer.CreatedAt = time.Now().Unix()
	answer.AnsweredBy = uniqueAlias
	answer.LikeNumber = 0
	// force MongoDB to always set its own generated ObjectIDs
	answer.Id = ""

	//sanitize answer
	html := utils.SanitizeHtml(answer.Answer)
	answer.Answer = html

	// insert the record
	insertResult, err := database.MG.Db.Collection(collection).InsertOne(c.Context(), answer)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	createdAnswer := new(Answer)

	// get the just inserted record in order to return it as response
	query = bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	database.MG.Db.Collection(collection).FindOne(c.Context(), query).Decode(&createdAnswer)

	//increase the totalAnswers by userId for user
	query = bson.D{{Key: "_id", Value: id}}
	updateQuery := bson.D{{Key: "$inc", Value: bson.D{{Key: "totalAnswers", Value: 1}}}}
	_, err = database.MG.Db.Collection("users").UpdateOne(c.Context(), query, updateQuery)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	//Streaks Logic
	currStreakTime := time.Now().Unix()

	type data struct {
		Streak         int64 `json:"streak" bson:"streak"`
		PrevStreakTime int64 `json:"prevStreakTime,omitempty" bson:"prevStreakTime,omitempty"`
	}
	var info data

	opts := options.FindOne().SetProjection(bson.D{{Key: "streak", Value: 1}, {Key: "prevStreakTime", Value: 1}})
	query = bson.D{{Key: "_id", Value: id}}
	err = database.MG.Db.Collection("users").FindOne(c.Context(), query, opts).Decode(&info)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	if info.Streak == 0 {
		info.Streak = 1
		info.PrevStreakTime = currStreakTime

		//update value of streak with userId in users
		query = bson.D{{Key: "_id", Value: id}}
		updateQuery = bson.D{{Key: "$set", Value: bson.D{{Key: "streak", Value: info.Streak}, {Key: "prevStreakTime", Value: info.PrevStreakTime}}}}
		_, err = database.MG.Db.Collection("users").UpdateOne(c.Context(), query, updateQuery)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
	}

	var limit24 int64 = 86400
	var limit48 int64 = 172800

	if (currStreakTime - info.PrevStreakTime) < limit24 {

		return c.Status(201).JSON(&fiber.Map{"streak": info.Streak, "userAnswer": createdAnswer})

	} else if (currStreakTime - info.PrevStreakTime) < limit48 {
		info.PrevStreakTime = currStreakTime
		info.Streak = info.Streak + 1
	}
	//update value of streak with userId in users
	query = bson.D{{Key: "_id", Value: id}}
	updateQuery = bson.D{{Key: "$set", Value: bson.D{{Key: "streak", Value: info.Streak}, {Key: "prevStreakTime", Value: info.PrevStreakTime}}}}
	_, err = database.MG.Db.Collection("users").UpdateOne(c.Context(), query, updateQuery)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(201).JSON(&fiber.Map{"streak": info.Streak, "userAnswer": createdAnswer})

}

func LikeAnwer(c *fiber.Ctx) error {

	answerId := c.Params("answerId")
	// user := c.Query("user")
	answeredBy := c.Query("answeredBy")
	numberId := c.Query("number")

	answeredBy = answeredBy + "#" + numberId

	//userId will taken from body
	userId := c.Locals("userId").(string)
	//put answerId into db with userId
	var like Like
	like.UserId = userId
	like.AnswerId = answerId
	_, err := database.MG.Db.Collection("likes").InsertOne(c.Context(), like)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	//increse likeNumber with answerId
	ansId, _ := primitive.ObjectIDFromHex(answerId)
	query := bson.D{{Key: "_id", Value: ansId}}
	updateQuery := bson.D{{Key: "$inc", Value: bson.D{{Key: "likeNumber", Value: 1}}}}
	database.MG.Db.Collection("answers").UpdateOne(c.Context(), query, updateQuery)

	//cal score
	//update db
	//increase totalLikes with answeredBy in user model
	query = bson.D{{Key: "uniqueAlias", Value: answeredBy}}
	updateQuery = bson.D{{Key: "$inc", Value: bson.D{{Key: "totalLikes", Value: 1}}}}
	database.MG.Db.Collection("users").UpdateOne(c.Context(), query, updateQuery)

	//add to redis, increment score
	// database.Redis.Client.ZIncrBy(c.Context(), "leaderboard", float64(score.LikeScore), user)
	//get the new rank
	// userRank, err := database.Redis.Client.ZRevRank(c.Context(), "leaderboard", userId).Result()
	// if err == redis.Nil {
	// 	// this means we didn't got an data from redis and send user -1 to signal no data
	// 	userRank = -1
	// }
	//send see if rank <=100 then sse

	return c.SendStatus(200)
}

func CancelLike(c *fiber.Ctx) error {

	answerId := c.Params("answerId")
	answeredBy := c.Query("answeredBy")
	numberId := c.Query("number")

	answeredBy = answeredBy + "#" + numberId
	//userId will taken from body
	userId := c.Locals("userId").(string)
	//delete answerId from db with userId
	query := bson.D{{Key: "$and", Value: bson.A{bson.D{{Key: "userId", Value: userId}}, bson.D{{Key: "answerId", Value: answerId}}}}}
	_, err := database.MG.Db.Collection("likes").DeleteOne(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	//decrease likeNumber with answerId
	ansId, _ := primitive.ObjectIDFromHex(answerId)
	query = bson.D{{Key: "_id", Value: ansId}}
	updateQuery := bson.D{{Key: "$inc", Value: bson.D{{Key: "likeNumber", Value: -1}}}}
	database.MG.Db.Collection("answers").UpdateOne(c.Context(), query, updateQuery)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	//decrease totalLikes with answeredBy in user model
	query = bson.D{{Key: "uniqueAlias", Value: answeredBy}}
	updateQuery = bson.D{{Key: "$inc", Value: bson.D{{Key: "totalLikes", Value: -1}}}}
	database.MG.Db.Collection("users").UpdateOne(c.Context(), query, updateQuery)

	return c.SendStatus(200)
}
