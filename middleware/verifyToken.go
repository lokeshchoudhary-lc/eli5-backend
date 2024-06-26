package middleware

import (
	"eli5/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthVerify(c *fiber.Ctx) error {

	refreshTokenCookie := c.Cookies("token")

	if refreshTokenCookie == "" {
		return c.Status(401).SendString("Unauthorized")
	} else {
		userId, err := auth.VerifyRefreshToken(refreshTokenCookie)

		if err != nil {
			return c.Status(401).SendString("Unauthorized")
		}

		c.Locals("userId", userId)
		return c.Next()
	}

}
func CompanyAuthVerify(c *fiber.Ctx) error {

	refreshTokenCookie := c.Cookies("companyToken")

	if refreshTokenCookie == "" {
		return c.Status(401).SendString("Unauthorized")
	} else {
		userId, err := auth.VerifyRefreshToken(refreshTokenCookie)

		if err != nil {
			return c.Status(401).SendString("Unauthorized")
		}

		c.Locals("userId", userId)
		return c.Next()
	}

}

// func AuthVerify(c *fiber.Ctx) error {

// 	headers := c.GetReqHeaders()
// 	authHeader := headers["Authorization"]
// 	tokenSplit := strings.Split(authHeader, "Bearer ")
// 	token := tokenSplit[1]

// 	userId, err := auth.VerifyRefreshToken(token)

// 	if err != nil {
// 		return c.Status(401).SendString("Unauthorized")
// 	}

// 	c.Locals("userId", userId)
// 	return c.Next()

// }

// accessTokenCookie := c.Cookies("accessToken")
// refreshTokenCookie := c.Cookies("refreshToken")

// if accessTokenCookie == "" {
// 	//
// 	if refreshTokenCookie == "" {
// 		// on 401 redirect to home page for login again
// 		return c.Status(401).SendString("Unauthorized")
// 	}

// 	userId, err := auth.VerifyRefreshToken(refreshTokenCookie)
// 	if err != nil {
// 		return c.Status(401).SendString("Unauthorized")
// 	}

// 	accessToken, err := auth.CreateAccessToken(userId)
// 	if err != nil {
// 		return c.Status(500).SendString(err.Error())
// 	}
// 	accessTokenCookie := fiber.Cookie{
// 		Name:     "accessToken",
// 		Value:    accessToken,
// 		SameSite: "Strict",
// 		Secure:   true,
// 		// Expires:  time.Now().Add(time.Minute * 15),
// 		MaxAge:   60 * 15,
// 		HTTPOnly: true,
// 	}
// 	c.Cookie(&accessTokenCookie)

// 	c.Locals("userId", userId)
// 	return c.Next()

// } else {
// 	userId, err := auth.VerifyAccessToken(accessTokenCookie)
// 	if err != nil {
// 		return c.Status(401).SendString("Unauthorized")
// 	}

// 	c.Locals("userId", userId)
// 	return c.Next()
// }
