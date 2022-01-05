package utils

import (
	"forum/forum/internal/service/delivery"
	"forum/forum/internal/service/repository"
	"forum/forum/internal/service/usecase"
	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
	"net/http"
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
	config := pgx.ConnConfig{
		User:                 "docker",
		Database:             "docker",
		Password:             "docker",
		PreferSimpleProtocol: false,
	}
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}
	db, err := pgx.NewConnPool(connPoolConfig)
	if err != nil {
		return nil, err
	}

	forumR := repository.NewForumRepository(db)
	postR := repository.NewPostRepository(db)
	serviceR := repository.NewServiceRepository(db)
	threadR := repository.NewThreadRepository(db)
	userR := repository.NewUserRepository(db)

	forumUC := usecase.NewForumUseCase(forumR)
	postUC := usecase.NewPostUseCase(postR)
	serviceUC := usecase.NewServiceUseCase(serviceR)
	threadUC := usecase.NewThreadUseCase(threadR)
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
	port := "5050"
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Error("app err = ", err)
		return err
	}
	return nil
}
