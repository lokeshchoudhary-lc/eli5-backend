package user

import (
	"eli5/auth"
	"eli5/config/database"
	"fmt"
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

		accessToken, err := auth.CreateAccessToken(userCheck.Id, userCheck.UniqueAlias)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		accessTokenCookie := fiber.Cookie{
			Name:     "accessToken",
			Value:    accessToken,
			Domain:   "eli5.club",
			SameSite: "none",
			Secure:   true,
			// Expires:  time.Now().Add(time.Minute * 15),
			MaxAge:   60 * 15,
			HTTPOnly: true,
		}

		refreshToken, err := auth.CreateRefreshToken(userCheck.Id, userCheck.UniqueAlias)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		refreshTokenCookie := fiber.Cookie{
			Name:     "refreshToken",
			Value:    refreshToken,
			Domain:   "eli5.club",
			SameSite: "none",
			Secure:   true,
			// Expires:  time.Now().Add(time.Hour * 168),
			MaxAge:   60 * 60 * 24 * 7,
			HTTPOnly: true,
		}
		// loginState := fiber.Cookie{
		// 	Name:     "loginState",
		// 	Value:    "true",
		// 	Expires:  time.Now().Add(time.Hour * 168),
		// 	HTTPOnly: false,
		// }
		c.Cookie(&accessTokenCookie)
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
	uniqueAlias := user.UniqueAlias

	// database.Redis.Client.ZAdd(c.Context(), "leaderboard", redis.Z{Score: float64(user.Score), Member: userId})

	accessToken, err := auth.CreateAccessToken(userId, uniqueAlias)
	if err != nil {

		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}
	accessTokenCookie := fiber.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Domain:   "eli5.club",
		SameSite: "none",
		Secure:   true,
		// Expires:  time.Now().Add(time.Minute * 15),
		MaxAge:   60 * 15,
		HTTPOnly: true,
	}

	refreshToken, err := auth.CreateRefreshToken(userId, uniqueAlias)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	refreshTokenCookie := fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Domain:   "eli5.club",
		SameSite: "none",
		Secure:   true,
		// Expires:  time.Now().Add(time.Hour * 168),
		MaxAge:   60 * 60 * 24 * 7,
		HTTPOnly: true,
	}
	// loginState := fiber.Cookie{
	// 	Name:     "loginState",
	// 	Value:    "true",
	// 	Expires:  time.Now().Add(time.Hour * 168),
	// 	HTTPOnly: false,
	// }
	// c.Cookie(&loginState)
	c.Cookie(&accessTokenCookie)
	c.Cookie(&refreshTokenCookie)

	return c.SendStatus(200)

}

func Logout(c *fiber.Ctx) error {

	accessTokenCookie := fiber.Cookie{
		Name:     "accessToken",
		Value:    "",
		Expires:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		HTTPOnly: true,
	}
	refreshTokenCookie := fiber.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Expires:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		HTTPOnly: true,
	}
	// loginState := fiber.Cookie{
	// 	Name:     "loginState",
	// 	Value:    "true",
	// 	Expires:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	// 	HTTPOnly: false,
	// }
	// c.Cookie(&loginState)
	c.Cookie(&accessTokenCookie)
	c.Cookie(&refreshTokenCookie)
	return c.SendStatus(200)
}
