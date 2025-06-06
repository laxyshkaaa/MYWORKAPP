package main

import (
	"alaricode/go-fiber/config"
	"alaricode/go-fiber/internal/home"
	"alaricode/go-fiber/internal/user"
	"alaricode/go-fiber/internal/vacancy"
	"alaricode/go-fiber/pkg/database"
	"alaricode/go-fiber/pkg/logger"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
)

func main() {
	config.Init()
	config.NewDatabaseConfig()
	logConfig := config.NewLogConfig()
	dbConfig := config.NewDatabaseConfig()
	customLogger := logger.NewLogger(logConfig)

	app := fiber.New()
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: customLogger,
	}))
	app.Use(recover.New())
	app.Static("/public", "./public")
	dbpool := database.CreateDbPool(dbConfig, customLogger)
	defer dbpool.Close()
	storage := postgres.New(postgres.Config{
		DB:         dbpool,
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	})
	store := session.New(session.Config{
		Storage: storage,
	})

	// Repositories
	vacancyRepo := vacancy.NewVacancyRepository(dbpool, customLogger)
	userRepo := user.NewUserRepository(dbpool)
	user.NewUserHandler(app, userRepo, store)

	// Handler
	home.NewHandler(app, customLogger, vacancyRepo, store)
	vacancy.NewHandler(app, customLogger, vacancyRepo, store)

	app.Listen(":3000")
}
