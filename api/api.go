package api

import (
	"eli5/api/answer"
	"eli5/api/company"
	"eli5/api/feed"
	"eli5/api/question"
	"eli5/api/user"
	"eli5/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupApiRoutes(app *fiber.App) {

	v1 := app.Group("/api/v1")

	//guest routes
	v1.Get("/leaderboard", user.GetLeaderBoard)
	v1.Get("/explore", question.Explore) // explore button (give all tags at once)
	v1.Get("/tags", question.GetTags)    // do pagiation on load more (query: page)
	v1.Get("/question/:questionId", question.GetQuestion)
	v1.Get("/tagStats/:tag", question.TagsPageStats)
	v1.Get("/questions/:tag", question.GetQuestionsOfTag)                        //do pagination on load more (after explore)(query: page)
	v1.Get("/changeQuestion/:questionId", question.ChangeQuestionWithPagination) //get the seleted single question for feed , have query: tag,page,action: next and back for question change
	v1.Get("/guestAnswers/:questionId", answer.GetGuestAnswers)                  //query : sort=latest,trending,page //end of page gives status 204
	v1.Get("/feed", feed.MakeHomeFeed)
	v1.Post("/completeProfile", user.CompleteProfile)
	v1.Get("/userCheck/:email", user.UserCheck)
	v1.Get("/gptAnswer/:questionId", answer.GetGptAnswer)
	// v1.Get("/userDetails/:username", user.GetProfileDetails)
	v1.Get("/profileDetails/:username", user.GetProfileDetails)
	v1.Post("/question/ask", question.AskQuestion)
	v1.Get("/authCheck", user.AuthCheck)

	//new landing page
	v1.Get("/allQuestions", question.GetAllQuestions)
	v1.Get("/trendingTags", question.GetTrendingTags)
	v1.Get("/backlinkQuestions/:tag", question.GetBacklinkQuestionsOfTag)

	//Comapny Routes
	v1.Get("/company", company.GetCompanyList)
	v1.Get("/company/:username", company.GetCompanyProfile)
	v1.Get("/companyQuestions/:companyId", company.GetCompanyQuestions)
	v1.Get("/companyCheck/:email", company.CompanyCheck)
	v1.Get("/companyAuthCheck", company.CompanyAuthCheck)
	v1.Get("/getCompanyUniqueAlias", middleware.CompanyAuthVerify, company.GetCompanyUniqueAlias)
	v1.Put("/company", middleware.CompanyAuthVerify, company.UpdateCompanyProfile)
	v1.Post("/company/like/:username", company.LikeCompany)
	v1.Post("/company/completeProfile", company.CompleteCompanyProfile)
	// v1.Get("/companyAnswer/:questionId", company.GetCompanyAnswer)
	// v1.Post("/companyAnswer/:questionId", middleware.CompanyAuthVerify, company.AnswerCompanyQuestion)

	//sse routes and login only
	// /leaderboard all broadcast
	// /details/:userId events: rank , likes? points?

	//action routes
	v1.Get("/answers/:questionId", middleware.AuthVerify, answer.GetAnswers)
	v1.Post("/answer/:questionId", middleware.AuthVerify, answer.PostAnswer)
	v1.Get("/userDetails", middleware.AuthVerify, user.GetUserDetails)
	v1.Put("/userDetails", middleware.AuthVerify, user.UpdateUserProfile)
	v1.Put("/like/:answerId", middleware.AuthVerify, answer.LikeAnwer)
	v1.Put("/cancelLike/:answerId", middleware.AuthVerify, answer.CancelLike)
	v1.Get("/userAnswer/:questionId", middleware.AuthVerify, answer.GetUserAnswer)
	v1.Get("/logout", user.Logout)
	v1.Get("/logout/company", company.Logout)

	//SSE routes
	// v1.Get("/sse/leaderboard", sse.LeaderboardSSE)
	// v1.Get("/sse/bc", sse.BroadcastSSE)

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

	// adminRoute.Post("/gptAnswer/:questionId", admin.PostGptAnswer)

	// adminRoute.Get("/tags", admin.GetTags)
	// adminRoute.Get("/questions/:tagName", admin.GetQuestionsOfTag)
	// adminRoute.Post("/choose/question/:questionId", admin.ChooseQuestion)
	// adminRoute.Post("/choose/tag/:tagName", admin.ChooseTag)
	// adminRoute.Post("/question", admin.PostQuestion)
	// adminRoute.Post("/tag/:tagName", admin.PostTag)
}
