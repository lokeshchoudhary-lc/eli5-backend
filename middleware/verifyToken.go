package middleware

import (
	"eli5/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthVerify(c *fiber.Ctx) error {
	accessTokenCookie := c.Cookies("accessToken")
	refreshTokenCookie := c.Cookies("refreshToken")

	if accessTokenCookie == "" {
		//
		if refreshTokenCookie == "" {
			// on 401 redirect to home page for login again
			return c.Status(401).SendString("Unauthorized")
		}

		userId, uniqueAlias, err := auth.VerifyRefreshToken(refreshTokenCookie)
		if err != nil {
			return c.Status(401).SendString("Unauthorized")
		}

		accessToken, err := auth.CreateAccessToken(userId, uniqueAlias)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		accessTokenCookie := fiber.Cookie{
			Name:     "accessToken",
			Value:    accessToken,
			SameSite: "none",
			Secure:   true,
			// Expires:  time.Now().Add(time.Minute * 15),
			MaxAge:   60 * 15,
			HTTPOnly: true,
		}
		c.Cookie(&accessTokenCookie)

		c.Locals("userId", userId)
		c.Locals("uniqueAlias", uniqueAlias)
		return c.Next()

	} else {
		userId, uniqueAlias, err := auth.VerifyAccessToken(accessTokenCookie)
		if err != nil {
			return c.Status(401).SendString("Unauthorized")
		}

		c.Locals("userId", userId)
		c.Locals("uniqueAlias", uniqueAlias)
		return c.Next()
	}

}
