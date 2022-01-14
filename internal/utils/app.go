package utils

import (
	"forum/internal/service/delivery"
	"forum/internal/service/repository"
	"forum/internal/service/usecase"
	database "forum/internal/utils/db"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	"github.com/jackc/pgx"
)

type App struct {
	forumDelivery   *delivery.ForumDelivery
	postDelivery    *delivery.PostDelivery
	serviceDelivery *delivery.ServiceDelivery
	threadDelivery  *delivery.ThreadDelivery
	userDelivery    *delivery.UserDelivery
	db              *pgx.ConnPool
}

func NewApp() (*App, error) {
	db, err := database.NewConnPool()
	if err != nil {
		return nil, err
	}
	err = database.Prepare(db)
	if err != nil {
		return nil, err
	}

	forumR := repository.NewForumRepository(db)
	postR := repository.NewPostRepository(db)
	serviceR := repository.NewServiceRepository(db)
	threadR := repository.NewThreadRepository(db)
	userR := repository.NewUserRepository(db)

	forumUC := usecase.NewForumUseCase(forumR, userR, threadR)
	postUC := usecase.NewPostUseCase(postR, userR, forumR, threadR)
	serviceUC := usecase.NewServiceUseCase(serviceR)
	threadUC := usecase.NewThreadUseCase(threadR, userR, postR)
	userUC := usecase.NewUserUseCase(userR)

	forumD := delivery.NewForumDelivery(forumUC)
	postD := delivery.NewPostDelivery(postUC)
	serviceD := delivery.NewServiceDelivery(serviceUC)
	threadD := delivery.NewThreadDelivery(threadUC)
	userD := delivery.NewUserDelivery(userUC)

	return &App{
		forumDelivery:   forumD,
		postDelivery:    postD,
		serviceDelivery: serviceD,
		threadDelivery:  threadD,
		userDelivery:    userD,
		db:              db,
	}, nil
}

func (a *App) Run() error {
	if a.db != nil {
		defer a.db.Close()
	}
	r := NewRouterForApp(a)
	log.Info("forum started")
	panic(fasthttp.ListenAndServe(":5000", r.HandleRequest))
	return nil
}
