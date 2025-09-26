package app

import (
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/pkg/logger"
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type App struct {
	Config   *config.Config
	Storage  *sql.DB
	Router   chi.Router
	Logger   logger.Logger
	Handlers handlers.Handlers
}

func New() *App {
	ctx := context.Background()

	logger, err := logger.NewLogger()
	if err != nil {
		log.Println("failed to create logger", err)
		panic(err)
	}
	logger.Info("logger initialized")

	config := config.MustGetConfig(ctx, logger)
	storage, err := sql.Open("sqlite3", config.SQLPath)
	if err != nil {
		logger.Error("failed to open database", zap.Error(err))
	}
	router := chi.NewRouter()

	handlers := initHandlers(logger, storage)

	app := App{
		Config:   config,
		Storage:  storage,
		Router:   router,
		Logger:   logger,
		Handlers: handlers,
	}

	app.addMiddleware()
	app.initRoutes()

	return &app
}

func initHandlers(logger logger.Logger, db *sql.DB) handlers.Handlers {
	return handlers.NewHandlers(logger)
}

func (app *App) Start() {
	addr := ":" + strconv.Itoa((app.Config.Port))
	app.Logger.Info("server started on port", addr)

	if err := http.ListenAndServe(addr, app.Router); err != nil {
		app.Logger.Error("failed to start server", zap.Error(err))
	}
}

func (app *App) initRoutes() {
	app.Router.Post("/echo", app.Handlers.Echo)
}
