package utils

import "github.com/gorilla/mux"

func NewRouterForApp(a *App) *mux.Router {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	//r.HandleFunc("/forum/create", a.forumDelivery.CreateForum)
	return r
}
