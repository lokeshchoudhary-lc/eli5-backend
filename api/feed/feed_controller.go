package feed

import (
	"eli5/api/question"
	"eli5/api/user"
	"eli5/config/database"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func MakeHomeFeed(c *fiber.Ctx) error {
	var userId string = c.Locals("userId").(string)
	var user user.User

	id, _ := primitive.ObjectIDFromHex(userId)
	query := bson.D{{Key: "_id", Value: id}}
	err := database.MG.Db.Collection("users").FindOne(c.Context(), query).Decode(&user)

	if err != nil {

		return c.Status(500).SendString(err.Error())
	}

	currStreakTime := time.Now().Unix()

	var limit48 int64 = 172800

	if user.PrevStreakTime != 0 {

		if (currStreakTime - user.PrevStreakTime) > limit48 {
			user.Streak = 0
			//update new value of streak into user with userId
			query := bson.D{{Key: "_id", Value: id}}
			updateQuery := bson.D{{Key: "$set", Value: bson.D{{Key: "streak", Value: user.Streak}}}}
			_, err := database.MG.Db.Collection("users").UpdateOne(c.Context(), query, updateQuery)
			if err != nil {
				return c.Status(500).SendString(err.Error())
			}
		}
	}

	//find all choosen tags
	query = bson.D{{Key: "choosen", Value: true}}

	cursor, err := database.MG.Db.Collection("tags").Find(c.Context(), query)

	if err != nil {
		fmt.Println("2")

		return c.Status(500).SendString(err.Error())
	}

	var tags []question.Tag = make([]question.Tag, 0)

	if err := cursor.All(c.Context(), &tags); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	// take each tag and find corresponding choosen question

	type feed struct {
		Tag      string            `json:"tag"`
		Question question.Question `json:"question"`
	}

	var data []feed = make([]feed, 0)

	for _, tag := range tags {

		var feed feed
		query = bson.D{{Key: "choosen", Value: true}, {Key: "tag", Value: tag.Tag}}
		feed.Tag = tag.Tag
		err := database.MG.Db.Collection("questions").FindOne(c.Context(), query).Decode(&feed.Question)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(500).SendString(err.Error())
			}
			return c.Status(500).SendString(err.Error())
		}
		data = append(data, feed)

	}

	if len(data) == 0 {
		return c.Status(500).SendString("No Data Found")
	}
	return c.JSON(&fiber.Map{
		"firstName":          user.FirstName,
		"profilePictureCode": user.ProfilePictureCode,
		"streak":             user.Streak,
		"totalLikes":         user.TotalLikes,
		"totalAnswers":       user.TotalAnswers,
		"feed":               data,
	})
}
