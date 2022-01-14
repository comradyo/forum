package utils

import (
	routing "github.com/qiangxue/fasthttp-routing"
)

func NewRouterForApp(a *App) *routing.Router {
	r := routing.New()
	r.Post("/api/forum/create", a.forumDelivery.CreateForum)
	r.Get("/api/forum/<slug>/details", a.forumDelivery.GetForumDetails)
	r.Post("/api/forum/<slug>/create", a.forumDelivery.CreateForumThread)
	r.Get("/api/forum/<slug>/users", a.forumDelivery.GetForumUsers)
	r.Get("/api/forum/<slug>/threads", a.forumDelivery.GetForumThreads)
	r.Get("/api/post/<id>/details", a.postDelivery.GetPostDetails)
	r.Post("/api/post/<id>/details", a.postDelivery.UpdatePostDetails)
	r.Post("/api/service/clear", a.serviceDelivery.Clear)
	r.Get("/api/service/status", a.serviceDelivery.GetStatus)
	r.Post("/api/thread/<slug_or_id>/create", a.threadDelivery.CreateThreadPosts)
	r.Get("/api/thread/<slug_or_id>/details", a.threadDelivery.GetThreadDetails)
	r.Post("/api/thread/<slug_or_id>/details", a.threadDelivery.UpdateThreadDetails)
	r.Get("/api/thread/<slug_or_id>/posts", a.threadDelivery.GetThreadPosts)
	r.Post("/api/thread/<slug_or_id>/vote", a.threadDelivery.VoteForThread)

	r.Post("/api/user/<nickname>/create", a.userDelivery.CreateUser)
	r.Get("/api/user/<nickname>/profile", a.userDelivery.GetUserProfile)
	r.Post("/api/user/<nickname>/profile", a.userDelivery.UpdateUserProfile)
	return r
}
