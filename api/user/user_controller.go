package user

import (
	"eli5/auth"
	"eli5/config/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetLeaderBoard(c *fiber.Ctx) error {

	var data []User = make([]User, 0)

	opts := options.Find().SetProjection(bson.D{{Key: "totalAnswers", Value: 1}, {Key: "totalLikes", Value: 1}, {Key: "uniqueAlias", Value: 1}, {Key: "profilePictureCode", Value: 1}}).SetSort(bson.D{{Key: "totalLikes", Value: -1}, {Key: "totalAnswers", Value: -1}}).SetLimit(100)
	query := bson.D{{}}
	cursor, err := database.MG.Db.Collection("users").Find(c.Context(), query, opts)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if err := cursor.All(c.Context(), &data); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(data)
}

func UserCheck(c *fiber.Ctx) error {

	userCheck := new(User)
	email := c.Params("email")

	query := bson.D{{Key: "email", Value: email}}
	err := database.MG.Db.Collection("users").FindOne(c.Context(), query).Decode(&userCheck)

	if err != mongo.ErrNoDocuments {
		refreshToken, err := auth.CreateRefreshToken(userCheck.Id)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		refreshTokenCookie := fiber.Cookie{
			Name:     "Token",
			Value:    refreshToken,
			SameSite: "Lax",
			Secure:   true,
			// Expires:  time.Now().Add(time.Hour * 168),
			MaxAge:   60 * 60 * 24 * 7,
			HTTPOnly: true,
		}
		c.Cookie(&refreshTokenCookie)
		// c.Cookie(&loginState)

		return c.Status(200).SendString("go_to_feed")
	}
	return c.Status(200).SendString("go_to_completeProfile")
}

func CompleteProfile(c *fiber.Ctx) error {
	user := new(User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	query := bson.D{{Key: "uniqueAlias", Value: user.UniqueAlias}}
	res := database.MG.Db.Collection("users").FindOne(c.Context(), query)
	if res.Err() != mongo.ErrNoDocuments {
		// uniqueAlias exists already
		return c.SendStatus(400)
	}

	user.CreatedAt = time.Now().Unix()
	user.Id = ""
	// user.Score = 0
	user.Streak = 0
	user.TotalAnswers = 0
	user.TotalLikes = 0

	insertedResult, err := database.MG.Db.Collection("users").InsertOne(c.Context(), user)
	if err != nil {
		c.Status(500).SendString(err.Error())
	}

	userId := insertedResult.InsertedID.(primitive.ObjectID).Hex()

	// database.Redis.Client.ZAdd(c.Context(), "leaderboard", redis.Z{Score: float64(user.Score), Member: userId}

	refreshToken, err := auth.CreateRefreshToken(userId)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	refreshTokenCookie := fiber.Cookie{
		Name:     "Token",
		Value:    refreshToken,
		SameSite: "Lax",
		Secure:   true,
		// Expires:  time.Now().Add(time.Hour * 168),
		MaxAge:   60 * 60 * 24 * 7,
		HTTPOnly: true,
	}

	c.Cookie(&refreshTokenCookie)

	return c.SendStatus(200)

}

func UpdateUserProfile(c *fiber.Ctx) error {

	var userId string = c.Locals("userId").(string)
	id, _ := primitive.ObjectIDFromHex(userId)

	type user struct {
		UniqueAlias  string `json:"uniqueAlias,omitempty" bson:"uniqueAlias,omitempty"`
		Bio          string `json:"bio,omitempty" bson:"bio,omitempty"`
		TwitterUrl   string `json:"twitterUrl,omitempty" bson:"twitterUrl,omitempty"`
		LinkedinUrl  string `json:"linkedinUrl,omitempty" bson:"linkedinUrl,omitempty"`
		InstagramUrl string `json:"instagramUrl,omitempty" bson:"instagramUrl,omitempty"`
	}

	userData := new(user)

	if err := c.BodyParser(userData); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	//check if new uniqueAlias/username already exists
	if userData.UniqueAlias != "" {
		query := bson.D{{Key: "uniqueAlias", Value: userData.UniqueAlias}}
		res := database.MG.Db.Collection("users").FindOne(c.Context(), query)
		if res.Err() != mongo.ErrNoDocuments {
			// uniqueAlias exists already
			return c.Status(400).SendString("exist")
		}
	}

	query := bson.D{{Key: "_id", Value: id}}
	updateQuery := bson.D{{Key: "$set", Value: userData}}
	_, err := database.MG.Db.Collection("users").UpdateOne(c.Context(), query, updateQuery)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(200)
}

func GetUserDetails(c *fiber.Ctx) error {

	userId := c.Locals("userId").(string)

	id, _ := primitive.ObjectIDFromHex(userId)

	var user User

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
			//update new value of streak  into user with id
			query := bson.D{{Key: "_id", Value: id}}
			updateQuery := bson.D{{Key: "$set", Value: bson.D{{Key: "streak", Value: user.Streak}}}}
			_, err := database.MG.Db.Collection("users").UpdateOne(c.Context(), query, updateQuery)
			if err != nil {
				return c.Status(500).SendString(err.Error())
			}
		}
	}

	// userRank, err := database.Redis.Client.ZRevRank(c.Context(), "leaderboard", userId).Result()
	// if err == redis.Nil {
	// 	// this means we didn't got an data from redis and send user -1 to signal no data
	// 	userRank = -1
	// }
	return c.JSON(&fiber.Map{
		"userId":             user.Id,
		"firstName":          user.FirstName,
		"uniqueAlias":        user.UniqueAlias,
		"profilePictureCode": user.ProfilePictureCode,
		"streak":             user.Streak,
		"totalLikes":         user.TotalLikes,
		"totalAnswers":       user.TotalAnswers,

		// "rank":               userRank,
	})
}

func GetProfileDetails(c *fiber.Ctx) error {
	username := c.Params("username")

	var user User

	query := bson.D{{Key: "uniqueAlias", Value: username}}
	err := database.MG.Db.Collection("users").FindOne(c.Context(), query).Decode(&user)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	currStreakTime := time.Now().Unix()

	var limit48 int64 = 172800

	if user.PrevStreakTime != 0 {

		if (currStreakTime - user.PrevStreakTime) > limit48 {
			user.Streak = 0
			//update new value of streak  into user with id
			query := bson.D{{Key: "uniqueAlias", Value: username}}
			updateQuery := bson.D{{Key: "$set", Value: bson.D{{Key: "streak", Value: user.Streak}}}}
			_, err := database.MG.Db.Collection("users").UpdateOne(c.Context(), query, updateQuery)
			if err != nil {
				return c.Status(500).SendString(err.Error())
			}
		}
	}

	// userRank, err := database.Redis.Client.ZRevRank(c.Context(), "leaderboard", userId).Result()
	// if err == redis.Nil {
	// 	// this means we didn't got an data from redis and send user -1 to signal no data
	// 	userRank = -1
	// }
	return c.JSON(&fiber.Map{
		"userId":             user.Id,
		"firstName":          user.FirstName,
		"bio":                user.Bio,
		"uniqueAlias":        user.UniqueAlias,
		"profilePictureCode": user.ProfilePictureCode,
		"streak":             user.Streak,
		"totalLikes":         user.TotalLikes,
		"totalAnswers":       user.TotalAnswers,
		"createdAt":          user.CreatedAt,
		"twitterUrl":         user.TwitterUrl,
		"instagramUrl":       user.InstagramUrl,
		"linkedinUrl":        user.LinkedinUrl,

		// "rank":               userRank,
	})
}

func Logout(c *fiber.Ctx) error {

	// accessTokenCookie := fiber.Cookie{
	// 	Name:     "accessToken",
	// 	Value:    "",
	// 	Expires:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	// 	HTTPOnly: true,
	// }
	refreshTokenCookie := fiber.Cookie{
		Name:     "Token",
		Value:    "",
		Expires:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		HTTPOnly: true,
	}
	// appStateCookie := fiber.Cookie{
	// 	Name:     "appState",
	// 	Value:    "",
	// 	Expires:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	// 	HTTPOnly: true,
	// }
	// loginState := fiber.Cookie{
	// 	Name:     "loginState",
	// 	Value:    "true",
	// 	Expires:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	// 	HTTPOnly: false,
	// }
	// c.Cookie(&loginState)

	// c.Cookie(&appStateCookie)
	// c.Cookie(&accessTokenCookie)
	c.Cookie(&refreshTokenCookie)
	return c.SendStatus(200)
}
