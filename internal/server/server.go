package server

import (
	"authstore/internal/config"
	"authstore/pkg/closer"
	"authstore/pkg/logging"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
)

func closerHandler(ctx context.Context) error {
	logger := logging.GetLogger()
	logger.Info("server closer handler")
	return nil
}
func Run(router *httprouter.Router, cfg *config.Config, shutup *closer.Shutdown) {
	logger := logging.GetLogger()
	shutup.Bind(closerHandler)
	var listener net.Listener
	var listenerErr error

	if cfg.Server.Listen.Type == "sock" {
		logger.Info("Listen UNIX socket")
		socketPath := getSocketPath()
		listener, listenerErr = net.Listen("unix", socketPath)
		logger.Infof("Server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("Listen TCP")
		listener, listenerErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Server.Listen.BindIP, cfg.Server.Listen.Port))
		logger.Infof("Server is listening port: %s:%s", cfg.Server.Listen.BindIP, cfg.Server.Listen.Port)
	}
	if listenerErr != nil {
		logger.Fatal(listenerErr)
	}
	server := &http.Server{
		Handler:      router,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Server.ReedTimeout) * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}

func getSocketPath() string {
	logger := logging.GetLogger()
	workingDir, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}
	appDir, err := filepath.Abs(filepath.Dir(workingDir))
	logger.Debug(appDir)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Create socket ")
	return path.Join(appDir, "web/build/app.sock")
}
