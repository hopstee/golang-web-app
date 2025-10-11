package app

import (
	"context"
	"database/sql"
	"log/slog"
	"mobile-backend-boilerplate/internal/config"
	"mobile-backend-boilerplate/internal/infrastructure"
	httpTransport "mobile-backend-boilerplate/internal/transport/http"
	"mobile-backend-boilerplate/internal/transport/http/options"
	"mobile-backend-boilerplate/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type App struct {
	Router *chi.Mux
	Config *config.Config
	DB     *sql.DB
	Logger *slog.Logger
}

func Init() (*App, error) {
	log := logger.New(slog.LevelDebug)

	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		log.Error("failed to load config", slog.Any("err", err))
		return nil, err
	}

	db, err := infrastructure.InitDB(cfg, log)
	if err != nil {
		log.Error("failed to initialize database", slog.Any("err", err))
		return nil, err
	}

	deps := infrastructure.Dependencies{
		DB:     db,
		Logger: log,
		Config: cfg,
	}

	if err := deps.InitRepos(); err != nil {
		return nil, err
	}

	deps.InitNotifiers()
	deps.InitServices()
	deps.InitHandlers()

	opts := options.Options{
		// constants
		StaticDir: cfg.Static.Dir,

		// services
		MobileAuthService: deps.MobileAuthService,
		AdminAuthService:  deps.AdminAuthService,
		PostService:       deps.PostService,

		// api handlers
		MobileAuthHandler: deps.MobileAuthHandler,
		AdminAuthHandler:  deps.AdminAuthHandler,
		RequestHandler:    deps.RequestHandler,

		// web handlers
		PostHandler:    deps.PostHandler,
		ContactHandler: deps.ContactHandler,

		// other options
		Logger: log,
	}

	router := httpTransport.NewRouter(opts)

	return &App{
		Router: router,
		Config: cfg,
		DB:     db,
		Logger: log,
	}, nil
}

func (a *App) Run() error {
	srv := &http.Server{
		Addr:    a.Config.Server.Addr,
		Handler: a.Router,
	}

	errCh := make(chan error, 1)
	go func() {
		a.Logger.Info("server starting", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			a.Logger.Error("server closed", slog.Any("err", err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stop:
		a.Logger.Info("shutting down server...")
	case err := <-errCh:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		a.Logger.Error("failed to shutdown server", slog.Any("err", err))
		return err
	}

	return nil
}

func (a *App) Close() error {
	if a.DB != nil {
		err := a.DB.Close()
		a.Logger.Info("Database connection closed")
		return err
	}
	return nil
}
