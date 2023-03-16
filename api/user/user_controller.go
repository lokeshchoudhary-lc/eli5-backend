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
			Name:  "accessToken",
			Value: accessToken,
			// Expires:  time.Now().Add(time.Minute * 15),
			MaxAge:   60 * 15,
			HTTPOnly: true,
		}

		refreshToken, err := auth.CreateRefreshToken(userCheck.Id, userCheck.UniqueAlias)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		refreshTokenCookie := fiber.Cookie{
			Name:  "refreshToken",
			Value: refreshToken,
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
	user.Rank = 0
	user.Streak = 0
	user.TotalAnswers = 0
	user.TotalLikes = 0

	insertedResult, err := database.MG.Db.Collection("users").InsertOne(c.Context(), user)
	if err != nil {
		c.Status(500).SendString(err.Error())
	}

	userId := insertedResult.InsertedID.(primitive.ObjectID).Hex()
	uniqueAlias := user.UniqueAlias

	accessToken, err := auth.CreateAccessToken(userId, uniqueAlias)
	if err != nil {

		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}
	accessTokenCookie := fiber.Cookie{
		Name:  "accessToken",
		Value: accessToken,
		// Expires:  time.Now().Add(time.Minute * 15),
		MaxAge:   60 * 15,
		HTTPOnly: true,
	}

	refreshToken, err := auth.CreateRefreshToken(userId, uniqueAlias)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	refreshTokenCookie := fiber.Cookie{
		Name:  "refreshToken",
		Value: refreshToken,
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

// func GoogleAuth(c *fiber.Ctx) error {
// 	path := auth.ConfigGoogle()
// 	url := path.AuthCodeURL("state")
// 	return c.Redirect(url)

// }

// func GoogleAuthCallback(c *fiber.Ctx) error {
// 	token, error := auth.ConfigGoogle().Exchange(c.Context(), c.FormValue("code"))
// 	if error != nil {
// 		return c.Status(500).SendString(error.Error())
// 	}

// 	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
// 	if err != nil {
// 		return c.Status(500).SendString(error.Error())
// 	}

// 	defer response.Body.Close()
// 	contents, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		return c.Status(500).SendString(err.Error())
// 	}
// 	type GoogleResponse struct {
// 		Email string `json:"email" `
// 		Name  string `json:"name" `
// 	}

// 	var data GoogleResponse
// 	errorz := json.Unmarshal(contents, &data)
// 	if errorz != nil {
// 		return c.Status(500).SendString(errorz.Error())
// 	}

// 	words := strings.Fields(data.Name)
// 	data.Name = words[0]

// 	userCheck := new(User)
// 	query := bson.D{{Key: "email", Value: data.Email}}
// 	err = database.MG.Db.Collection("users").FindOne(c.Context(), query).Decode(&userCheck)

// 	if err != nil {

// 		accessToken, err := auth.CreateAccessToken(userCheck.Id, userCheck.UniqueAlias)
// 		if err != nil {
// 			return c.Status(500).SendString(err.Error())
// 		}
// 		accessTokenCookie := fiber.Cookie{
// 			Name:     "accessToken",
// 			Value:    accessToken,
// 			Expires:  time.Now().Add(time.Minute * 15),
// 			HTTPOnly: true,
// 		}

// 		refreshToken, err := auth.CreateRefreshToken(userCheck.Id, userCheck.UniqueAlias)
// 		if err != nil {
// 			return c.Status(500).SendString(err.Error())
// 		}
// 		refreshTokenCookie := fiber.Cookie{
// 			Name:     "refreshToken",
// 			Value:    refreshToken,
// 			Expires:  time.Now().Add(time.Hour * 168),
// 			HTTPOnly: true,
// 		}
// 		c.Cookie(&accessTokenCookie)
// 		c.Cookie(&refreshTokenCookie)

// 		return c.Status(200).SendString("go_to_feed")
// 	}

// 	return c.Status(200).JSON(data)
// }
