package company

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

func CompanyAuthCheck(c *fiber.Ctx) error {

	refreshTokenCookie := c.Cookies("companyToken")

	if refreshTokenCookie == "" {
		return c.JSON(&fiber.Map{
			"auth": false,
		})
	} else {
		return c.JSON(&fiber.Map{
			"auth": true,
		})
	}

}
func CompanyCheck(c *fiber.Ctx) error {

	companyCheck := new(Company)
	email := c.Params("email")

	query := bson.D{{Key: "email", Value: email}}
	err := database.MG.Db.Collection("company").FindOne(c.Context(), query).Decode(&companyCheck)

	if err != mongo.ErrNoDocuments {
		refreshToken, err := auth.CreateRefreshToken(companyCheck.Id)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		refreshTokenCookie := fiber.Cookie{
			Name:     "companyToken",
			Value:    refreshToken,
			SameSite: "Lax",
			Secure:   true,
			// Expires:  time.Now().Add(time.Hour * 168),
			MaxAge:   60 * 60 * 24 * 7,
			HTTPOnly: true,
		}
		c.Cookie(&refreshTokenCookie)
		// c.Cookie(&loginState)

		return c.Status(200).SendString("go_to_profile")
	}
	return c.Status(200).SendString("go_to_completeProfile")
}

func GetCompanyList(c *fiber.Ctx) error {
	//give list of all companies present
	options := options.Find().SetProjection(bson.D{{Key: "uniqueAlias", Value: 1}, {Key: "bio", Value: 1}, {Key: "name", Value: 1}, {Key: "brandColor", Value: 1}, {Key: "profilePictureUrl", Value: 1}})
	query := bson.D{{}}
	cursor, err := database.MG.Db.Collection("company").Find(c.Context(), query, options)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	var allCompany []Company = make([]Company, 0)

	if err := cursor.All(c.Context(), &allCompany); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(allCompany)
}

func GetCompanyProfile(c *fiber.Ctx) error {
	username := c.Params("username")

	company := new(Company)

	// opts := options.FindOne().SetProjection(bson.D{{Key: "email", Value: -1}})
	query := bson.D{{Key: "uniqueAlias", Value: username}}
	err := database.MG.Db.Collection("company").FindOne(c.Context(), query).Decode(&company)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	company.Email = ""

	return c.JSON(company)

}

func GetCompanyQuestions(c *fiber.Ctx) error {
	companyId := c.Params("companyId")

	var companyQuestions []CompanyQuestion = make([]CompanyQuestion, 0)

	query := bson.D{{Key: "companyId", Value: companyId}, {Key: "choosen", Value: true}}
	cursor, err := database.MG.Db.Collection("companyQuestions").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if err := cursor.All(c.Context(), &companyQuestions); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if length := len(companyQuestions); length == 0 {
		return c.SendStatus(204)
	}

	return c.JSON(companyQuestions)

}
func CompleteCompanyProfile(c *fiber.Ctx) error {
	company := new(Company)

	if err := c.BodyParser(company); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	query := bson.D{{Key: "uniqueAlias", Value: company.UniqueAlias}}
	res := database.MG.Db.Collection("company").FindOne(c.Context(), query)
	if res.Err() != mongo.ErrNoDocuments {
		// uniqueAlias exists already
		return c.SendStatus(400)
	}

	company.Id = ""
	company.TotalLikes = 0
	company.CreatedAt = time.Now().Unix()

	insertedResult, err := database.MG.Db.Collection("company").InsertOne(c.Context(), company)
	if err != nil {
		c.Status(500).SendString(err.Error())
	}

	companyId := insertedResult.InsertedID.(primitive.ObjectID).Hex()

	refreshToken, err := auth.CreateRefreshToken(companyId)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	refreshTokenCookie := fiber.Cookie{
		Name:     "companyToken",
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

func LikeCompany(c *fiber.Ctx) error {

	username := c.Params("username")

	//increse totalLikes with username
	query := bson.D{{Key: "uniqueAlias", Value: username}}
	updateQuery := bson.D{{Key: "$inc", Value: bson.D{{Key: "totalLikes", Value: 1}}}}
	database.MG.Db.Collection("company").UpdateOne(c.Context(), query, updateQuery)

	return c.SendStatus(200)

}

func GetCompanyUniqueAlias(c *fiber.Ctx) error {
	var userId string = c.Locals("userId").(string)
	id, _ := primitive.ObjectIDFromHex(userId)

	company := new(Company)

	opts := options.FindOne().SetProjection(bson.D{{Key: "uniqueAlias", Value: 1}})
	query := bson.D{{Key: "_id", Value: id}}
	err := database.MG.Db.Collection("company").FindOne(c.Context(), query, opts).Decode(&company)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(company)

}
func UpdateCompanyProfile(c *fiber.Ctx) error {

	var userId string = c.Locals("userId").(string)
	id, _ := primitive.ObjectIDFromHex(userId)

	type company struct {
		UniqueAlias       string `json:"uniqueAlias,omitempty" bson:"uniqueAlias,omitempty"`
		Name              string `json:"name,omitempty" bson:"name,omitempty"`
		Bio               string `json:"bio,omitempty" bson:"bio,omitempty"`
		BrandColor        string `json:"brandColor,omitempty" bson:"brandColor,omitempty"`
		WebsiteUrl        string `json:"websiteUrl,omitempty" bson:"websiteUrl,omitempty"`
		TwitterUrl        string `json:"twitterUrl,omitempty" bson:"twitterUrl,omitempty"`
		LinkedinUrl       string `json:"linkedinUrl,omitempty" bson:"linkedinUrl,omitempty"`
		FacebookUrl       string `json:"facebookUrl,omitempty" bson:"facebookUrl,omitempty"`
		YoutubeUrl        string `json:"youtubeUrl,omitempty" bson:"youtubeUrl,omitempty"`
		InstagramUrl      string `json:"instagramUrl,omitempty" bson:"instagramUrl,omitempty"`
		OtherUrl          string `json:"otherUrl,omitempty" bson:"otherUrl,omitempty"`
		ProfilePictureUrl string `json:"profilePictureUrl,omitempty" bson:"profilePictureUrl,omitempty"`
	}

	companyData := new(company)

	if err := c.BodyParser(companyData); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	//check if new uniqueAlias/username already exists
	if companyData.UniqueAlias != "" {
		query := bson.D{{Key: "uniqueAlias", Value: companyData.UniqueAlias}}
		res := database.MG.Db.Collection("company").FindOne(c.Context(), query)
		if res.Err() != mongo.ErrNoDocuments {
			// uniqueAlias exists already
			return c.Status(400).SendString("exist")
		}
	}

	query := bson.D{{Key: "_id", Value: id}}
	updateQuery := bson.D{{Key: "$set", Value: companyData}}
	_, err := database.MG.Db.Collection("company").UpdateOne(c.Context(), query, updateQuery)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(200)

}
