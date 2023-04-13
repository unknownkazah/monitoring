package main

import (
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"monitoring/health"
)

func main() {
	// Get os env
	googleURL := os.Getenv("GOOGLE_URL")
	libraryURL := os.Getenv("LIBRARY_URL")
	telegramToken := os.Getenv("TELEGRAM_TOKEN")

	// Init telegram bot
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// init handlers
	services := make([]health.Service, 0)

	services = append(services,
		health.Service{
			Name: "google",
			URL:  googleURL,
		})

	services = append(services,
		health.Service{
			Name: "library",
			URL:  libraryURL,
		})

	healthHandler := health.NewHandler(services)

	monitoringService := monitoring.NewService(bot)

	// init cron tab
	crontab := cron.New()
	crontab.AddFunc("@every 1h", func() {
		monitoringService(-1001626172300)
	})
	crontab.Start()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	apiGroup := e.Group("/api")

	// Authors
	healthGroup := apiGroup.Group("/health")
	healthGroup.GET("/:id", healthHandler.Check)
}
