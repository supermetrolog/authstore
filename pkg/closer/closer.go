package closer

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	exitOK       int8 = 0
	exitError    int8 = 1
	closeTimeout int  = 15
)

type handler func(ctx context.Context) error

type Logging interface {
	Info(msg ...any)
	Warn(msg ...any)
}

type Shutdown struct {
	handlers []func(ctx context.Context) error
	logger   Logging
}

func New(logger Logging) *Shutdown {
	return &Shutdown{
		logger: logger,
	}
}

func (c *Shutdown) Bind(h handler, handlers ...handler) {
	c.handlers = append(c.handlers, h)
	for _, h := range handlers {
		c.handlers = append(c.handlers, h)
	}
}

func (c *Shutdown) Close(ctx context.Context) []error {
	var handleError []error
	handlersQueue := make(chan handler, 1)
	go func() {
		// calling handlers in reverse order
		for i := len(c.handlers); i > 0; i-- {
			handlersQueue <- c.handlers[i-1]
		}
		close(handlersQueue)
	}()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("func Close ended with context timeout")
			return handleError
		case handl, ok := <-handlersQueue:
			if !ok {
				c.logger.Info("func Close ended with closed chan")
				return handleError
			}
			if err := handl(ctx); err != nil {
				handleError = append(handleError, err)
			}
		}
	}
}

func (c *Shutdown) Listen() {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(
		sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		s := <-sigChan
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(closeTimeout)*time.Second)
		defer cancel()
		c.logger.Info("gracefull shutdown ", s)
		var exitStatus int8 = exitOK
		if errs := c.Close(ctx); errs != nil {
			c.logger.Warn("closed handlers ended with a errors ", errs)
			exitStatus = exitError
		}
		c.logger.Info("exit with status: ", exitStatus)
		os.Exit(int(exitStatus))
	}()
}
