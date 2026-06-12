package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/yourusername/user-api/config"
	sqlcdb "github.com/yourusername/user-api/db/sqlc"
	"github.com/yourusername/user-api/internal/handler"
	"github.com/yourusername/user-api/internal/logger"
	"github.com/yourusername/user-api/internal/middleware"
	"github.com/yourusername/user-api/internal/repository"
	"github.com/yourusername/user-api/internal/routes"
	"github.com/yourusername/user-api/internal/service"
)

func main() {
	// ------------------------------------------------------------------
	// 1. Logger
	// ------------------------------------------------------------------
	if err := logger.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialise logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	log := logger.Get()

	// ------------------------------------------------------------------
	// 2. Config
	// ------------------------------------------------------------------
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	// ------------------------------------------------------------------
	// 3. Database
	// ------------------------------------------------------------------
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		log.Fatal("failed to open database connection", zap.Error(err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("database is unreachable", zap.Error(err))
	}
	log.Info("connected to MySQL", zap.String("db", cfg.DBName))

	// ------------------------------------------------------------------
	// 4. Wire dependencies
	// ------------------------------------------------------------------
	queries := sqlcdb.New(db)
	repo := repository.NewUserRepository(queries)
	svc := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(svc, log)

	// ------------------------------------------------------------------
	// 5. Fiber app
	// ------------------------------------------------------------------
	app := fiber.New(fiber.Config{
		// Return structured JSON for 404 / method-not-allowed built-ins.
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			msg := "internal server error"
			var e *fiber.Error
			if ferr, ok := err.(*fiber.Error); ok {
				e = ferr
				code = e.Code
				msg = e.Message
			}
			return c.Status(code).JSON(fiber.Map{"error": msg})
		},
	})

	// Global middleware
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger(log))

	// Routes
	routes.Register(app, userHandler)

	// ------------------------------------------------------------------
	// 6. Graceful shutdown
	// ------------------------------------------------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		addr := ":" + cfg.AppPort
		log.Info("server starting", zap.String("addr", addr))
		if err := app.Listen(addr); err != nil {
			log.Error("server error", zap.Error(err))
		}
	}()

	<-quit
	log.Info("shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Error("server shutdown error", zap.Error(err))
	}
	log.Info("server exited")
}
