package api

import (
	"authstore/internal/config"
	accesshttp "authstore/internal/domain/access/delivery/http"
	accessmysqlrepo "authstore/internal/domain/access/repository/mysql"
	accessService "authstore/internal/domain/access/service"
	treehttp "authstore/internal/domain/tree/delivery/http"
	treemysqlrepo "authstore/internal/domain/tree/repository/mysql"
	treeService "authstore/internal/domain/tree/service"
	userhttp "authstore/internal/domain/user/delivery/http"
	usermysqlrepo "authstore/internal/domain/user/repository/mysql"
	userservice "authstore/internal/domain/user/service"
	"authstore/internal/server"
	"authstore/pkg/client/mysql"
	"authstore/pkg/closer"
	"authstore/pkg/logging"
	"context"

	"github.com/julienschmidt/httprouter"
)

func closerHandler(ctx context.Context) error {
	logger := logging.GetLogger()
	logger.Info("api closer handler")
	return nil
}

type Handler interface {
	Register(router *httprouter.Router)
}

func registerHandlers(router *httprouter.Router, handlers ...Handler) {
	for _, handler := range handlers {
		handler.Register(router)
	}
}

func Run() {
	logger := logging.GetLogger()
	logger.Info("run API server")

	router := httprouter.New()

	cfg := config.GetConfig()

	shutdown := closer.New(logging.GetLogger())
	shutdown.Listen()
	shutdown.Bind(closerHandler)

	logger.Info("handlers register")
	client, err := mysql.NewClient()
	if err != nil {
		logger.Fatalf("mysql client error %v", err)
	}
	userService := userservice.NewService(
		logger,
		usermysqlrepo.NewRepository(logger, client),
		accessmysqlrepo.NewRepository(logger, client),
	)
	userHandler := userhttp.NewHandler(
		logger,
		userService,
	)

	accessHandler := accesshttp.NewHandler(
		logger,
		accessService.NewService(
			logger,
			accessmysqlrepo.NewRepository(logger, client),
		),
		userService,
	)
	treeHandler := treehttp.NewHandler(
		logger,
		treeService.NewService(
			logger,
			treemysqlrepo.NewRepository(logger, client),
		),
		userService,
	)
	registerHandlers(router,
		userHandler,
		accessHandler,
		treeHandler,
	)

	logger.Info("run server")
	server.Run(router, cfg, shutdown)
}
