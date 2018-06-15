package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndrewBurian/adulting-api/data"
	"github.com/AndrewBurian/adulting-api/middlewares"
	"github.com/AndrewBurian/powermux"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
)

const (
	// DefaultPort is the port used if none is specified
	DefaultPort = "8080"
)

func main() {

	debug := flag.Bool("v", false, "Debug verbosity")
	quiet := flag.Bool("q", false, "Errors only")
	flag.Parse()

	// Setup logging
	// ------------------------------------------------------------------------
	if _, found := os.LookupEnv("DEBUG"); found {
		*debug = true
	}

	if *debug && *quiet {
		log.Fatal("Can only set one of -q and -v")
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else if *quiet {
		log.SetLevel(log.ErrorLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Setup Database
	// ------------------------------------------------------------------------
	dbURL, found := os.LookupEnv("DATABASE_URL")
	if !found {
		log.Fatal("DATABASE_URL required")
	}

	dbOpts, err := pg.ParseURL(dbURL)
	if err != nil {
		log.WithError(err).Fatal("Could not parse DB connection string")
	}

	dbConn := pg.Connect(dbOpts)

	// Mocks
	// ------------------------------------------------------------------------
	mockUser := data.NewMockUserDal()
	mockActivity := data.NewMockActivityDal()

	// Setup Server
	// ------------------------------------------------------------------------
	mux := powermux.NewServeMux()
	mux.MiddlewareFunc("/", middlewares.ContentTypeDetect)

	authDetect := &middlewares.AuthDetection{
		DB: mockUser,
	}
	mux.Middleware("/", authDetect)

	port, found := os.LookupEnv("PORT")
	if !found {
		port = DefaultPort
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	// graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-signals
		log.WithField("signal", s).Info("Trapped signal")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		err := server.Shutdown(shutdownCtx)
		if err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("Error shutting down server")
		}
	}()

	// Setup Handlers
	// ------------------------------------------------------------------------
	users := &UserHandler{
		db: dbConn,
	}
	users.Setup(mux.Route("/user"))

	auth := &AuthHandler{
		db:     mockUser,
		logger: log.WithField("component", "auth"),
	}
	auth.Setup(mux.Route("/auth"))

	activity := &ActivityHandler{
		db:     mockActivity,
		logger: log.WithField("component", "activity"),
	}
	activity.Setup(mux.Route("/activity"))

	if *debug {
		fmt.Println(mux)
	}

	// Run
	// ------------------------------------------------------------------------
	log.WithField("addr", server.Addr).Info("Server Starting")
	if err = server.ListenAndServe(); err != http.ErrServerClosed {
		log.WithError(err).Fatal("Error with server")
	}

	log.Info("Server Shutting down")

}
