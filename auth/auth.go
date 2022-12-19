package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// ConfigGoogle to set config of oauth
func ConfigGoogle() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}
	return conf
}

var jwtAccessKey []byte = []byte(os.Getenv("JWT_ACCESS_KEY"))
var jwtRefreshKey []byte = []byte(os.Getenv("JWT_REFRESH_KEY"))

type MyCustomClaims struct {
	UserId      string `json:"userId"`
	UniqueAlias string `json:"uniqueAlias"`
	jwt.RegisteredClaims
}

func CreateAccessToken(userId string, uniqueAlias string) (string, error) {
	claims := MyCustomClaims{
		userId,
		uniqueAlias,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    "Eli5",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(jwtAccessKey)
	return tokenString, err
}

func CreateRefreshToken(userId string, uniqueAlias string) (string, error) {
	claims := MyCustomClaims{
		userId,
		uniqueAlias,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)),
			Issuer:    "Eli5",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(jwtRefreshKey)
	return tokenString, err

}
func VerifyAccessToken(tokenString string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtAccessKey, nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		userId := claims.UserId
		uniqueAlias := claims.UniqueAlias

		return userId, uniqueAlias, nil
	}

	return "", "", err

}
func VerifyRefreshToken(tokenString string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtRefreshKey, nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		userId := claims.UserId
		uniqueAlias := claims.UniqueAlias

		return userId, uniqueAlias, nil
	}

	return "", "", err
}
