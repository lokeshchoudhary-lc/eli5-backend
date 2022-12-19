package api

import (
	"eli5/api/answer"
	"eli5/api/feed"
	"eli5/api/user"
	"eli5/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupApiRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")
	v1.Get("/feed", middleware.AuthVerify, feed.MakeHomeFeed)

	v1.Get("/userAnswer/:questionId", middleware.AuthVerify, answer.GetUserAnswer)
	v1.Get("/answers/:questionId", middleware.AuthVerify, answer.GetAnswers) //query : sort=latest,trending,page //end of page gives status 204
	v1.Post("/answer/:questionId", middleware.AuthVerify, answer.PostAnswer)
	v1.Put("/like/:answerId", middleware.AuthVerify, answer.LikeAnwer)
	v1.Put("/cancelLike/:answerId", middleware.AuthVerify, answer.CancelLike)

	v1.Get("/leaderboard", middleware.AuthVerify, user.GetLeaderBoard)
	v1.Post("/completeProfile", user.CompleteProfile)
	v1.Get("/userCheck/:email", user.UserCheck)
	v1.Get("/logout", user.Logout)
	// v1.Get("/auth/google", user.GoogleAuth)
	// v1.Get("/auth/google/callback", user.GoogleAuthCallback)
	// v1.Post("/signup")

	// adminRoute := v1.Group("/admin")

	// adminRoute.Get("/tags", admin.GetTags)
	// adminRoute.Get("/questions/:tagName", admin.GetQuestionsOfTag)
	// adminRoute.Post("/choose/question/:questionId", admin.ChooseQuestion)
	// adminRoute.Post("/choose/tag/:tagName", admin.ChooseTag)
	// adminRoute.Post("/question", admin.PostQuestion)
	// adminRoute.Post("/tag/:tagName", admin.PostTag)
}
