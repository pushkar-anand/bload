package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pushkar-anand/bload/storage"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer  http.Server
	handler    http.Handler

	router     *mux.Router

	client     *redis.Client
	logger     *logrus.Logger

	killServer chan int
	connClose  chan int
}

func NewServer(logger *logrus.Logger, client *redis.Client) *Server {
	r := mux.NewRouter()
	r.StrictSlash(false)

	s := &Server{
		router:     r,
		logger:     logger,
		client:     client,
		killServer: make(chan int),
		connClose:  make(chan int),
	}

	s.initialize()

	return s
}

func (s *Server) initialize() {
	go func() {
		is := make(chan os.Signal, 1)
		signal.Notify(is, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

		select {
		case i := <-is:
			s.logger.Infof("Received System Interrupt: %v", i)
		case <-s.killServer:
			s.logger.Info("Received Kill Request")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		s.logger.Info("Stopped listening to incoming connections")
		s.logger.Info("Waiting for open connections to close")

		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			s.logger.WithError(err).Error("Error shutting down server")
		}

		close(s.connClose)
	}()

	s.addRoutes()
	s.addMiddleware()
	s.addHandlers()
}

func (s *Server) addRoutes() {
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	str := storage.NewStorage(s.client)
	sh := storage.NewHandler(str, s.logger)
	storage.AddRoutes(s.router, sh)
}

func (s *Server) addMiddleware() {
	//s.router.Use()
}

func (s *Server) addHandlers() {
	s.handler = s.router
	s.handler = handlers.CombinedLoggingHandler(os.Stdout, s.handler)
	s.handler = handlers.CompressHandler(s.handler)
}

func (s *Server) Listen() {
	addr := fmt.Sprintf("%s:%s", "", os.Getenv(envPort))
	s.httpServer = http.Server{
		Addr:         addr,
		Handler:      s.handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.logger.Infof("Server started at: %s", addr)

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.WithError(err).Panic("error starting server")
	}
}
