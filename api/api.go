package api

import (
	"eli5/api/answer"
	"eli5/api/feed"
	"eli5/api/question"
	"eli5/api/user"
	"eli5/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupApiRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	//prototype 3
	//guest routes
	v1.Get("/leaderboard", user.GetLeaderBoard)
	v1.Get("/explore", question.Explore) // explore button (give all tags at once)
	v1.Get("/tags", question.GetTags)    // do pagiation on load more (query: page)
	v1.Get("/question/:questionId", question.GetQuestion)
	v1.Get("/questions/:tag", question.GetQuestionsOfTag)                        //do pagination on load more (after explore)(query: page)
	v1.Get("/changeQuestion/:questionId", question.ChangeQuestionWithPagination) //get the seleted single question for feed , have query: tag,page,action: next and back for question change
	v1.Get("/answers/:questionId", answer.GetAnswers)                            //query : sort=latest,trending,page //end of page gives status 204
	v1.Get("/feed", feed.MakeHomeFeed)
	v1.Post("/completeProfile", user.CompleteProfile)
	v1.Get("/userCheck/:email", user.UserCheck)

	//sse routes and login only
	// /leaderboard all broadcast
	// /details/:userId events: rank , likes? points?

	//action routes
	v1.Post("/answer/:questionId", middleware.AuthVerify, answer.PostAnswer)
	v1.Put("/like/:answerId", middleware.AuthVerify, answer.LikeAnwer)
	v1.Put("/cancelLike/:answerId", middleware.AuthVerify, answer.CancelLike)
	v1.Post("/question/ask", middleware.AuthVerify, question.AskQuestion)
	v1.Get("/userDetails", middleware.AuthVerify, feed.GetUserDetails)
	v1.Get("/userAnswer/:questionId", middleware.AuthVerify, answer.GetUserAnswer)
	v1.Get("/logout", user.Logout)

	// v1.Get("/singleQuestion/:tag", question.GetQuestionOfTag)           //do pagination on load more (after explore)(query: page)
	// v1.Get("/feed", middleware.AuthVerify, feed.MakeHomeFeed)

	// v1.Post("/answer/:questionId", middleware.AuthVerify, answer.PostAnswer)
	// v1.Put("/like/:answerId", middleware.AuthVerify, answer.LikeAnwer)
	// v1.Put("/cancelLike/:answerId", middleware.AuthVerify, answer.CancelLike)

	// v1.Get("/leaderboard", middleware.AuthVerify, user.GetLeaderBoard)
	// v1.Post("/completeProfile", user.CompleteProfile)
	// v1.Get("/userCheck/:email", user.UserCheck)

	// prototype stage 2
	// v1.Get("/explore", middleware.AuthVerify, question.Explore)                                // explore button (give all tags at once)
	// v1.Get("/tags", middleware.AuthVerify, question.GetTags)                                   // do pagiation on load more (query: page)
	// v1.Get("/singleQuestion/:tag", middleware.AuthVerify, question.GetQuestionOfTag)           //do pagination on load more (after explore)(query: page)
	// v1.Get("/questions/:tag", middleware.AuthVerify, question.GetQuestionsOfTag)               //do pagination on load more (after explore)(query: page)
	// v1.Get("/question/:questionId", middleware.AuthVerify, question.GetQuestionWithPagination) //get the seleted single question for feed , have query: tag,page,action: next and back for question change

	////////////////////////////////////////////////////////////////////////////////

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
