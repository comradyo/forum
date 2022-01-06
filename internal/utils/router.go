package utils

import "github.com/gorilla/mux"

func NewRouterForApp(a *App) *mux.Router {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.HandleFunc("/forum/create", a.forumDelivery.CreateForum).Methods("POST")
	r.HandleFunc("/forum/{slug}/details", a.forumDelivery.GetForumDetails).Methods("GET")
	r.HandleFunc("/forum/{slug}/create", a.forumDelivery.CreateForumThread).Methods("POST")
	r.HandleFunc("/forum/{slug}/users", a.forumDelivery.GetForumUsers).Methods("GET")
	r.HandleFunc("/forum/{slug}/threads", a.forumDelivery.GetForumThreads).Methods("GET")
	r.HandleFunc("/post/{id}/details", a.postDelivery.GetPostDetails).Methods("GET")
	r.HandleFunc("/post/{id}/details", a.postDelivery.UpdatePostDetails).Methods("POST")
	r.HandleFunc("/service/clear", a.serviceDelivery.Clear).Methods("POST")
	r.HandleFunc("/service/status", a.serviceDelivery.GetStatus).Methods("GET")
	r.HandleFunc("/thread/{slug_or_id}/create", a.threadDelivery.CreateThreadPosts).Methods("POST")
	r.HandleFunc("/thread/{slug_or_id}/details", a.threadDelivery.GetThreadDetails).Methods("GET")
	r.HandleFunc("/thread/{slug_or_id}/details", a.threadDelivery.UpdateThreadDetails).Methods("POST")
	r.HandleFunc("/thread/{slug_or_id}/posts", a.threadDelivery.GetThreadPosts).Methods("GET")
	r.HandleFunc("/thread/{slug_or_id}/vote", a.threadDelivery.VoteForThread).Methods("POST")
	r.HandleFunc("/user/{nickname}/create", a.userDelivery.CreateUser).Methods("POST")
	r.HandleFunc("/user/{nickname}/profile", a.userDelivery.GetUserProfile).Methods("GET")
	r.HandleFunc("/user/{nickname}/profile", a.userDelivery.UpdateUserProfile).Methods("POST")
	return r
}
