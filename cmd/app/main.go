package main

import (
	"context"

	_ "library/docs"
	app "library/internal/app"
	"library/internal/config"
	"library/internal/http-server/handlers"
	"library/internal/lib/logger"
	"library/internal/lib/logger/sl"
	"library/internal/repositoy/postgres"
	"library/internal/service"
	"library/internal/storage/postgres"
)

//	@title			Library API
//	@version		1.0
//	@description	This API allows you to manage songs in the library system.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.example.com/support
//	@contact.email	support@example.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host	localhost:8000

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log.Info("start application")

	ctx := context.TODO()
	storage, err := postgres.New(ctx, cfg)
	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
		panic(err)
	}
	log.Info("successful connection to the database")
	defer storage.Close()
	repo := repository.NewStorage(storage, log)

	serv := service.NewService(log, repo)
	h := handlers.NewHandler(*serv, *log)
	router := app.SetupRouter(h, log)

	server := app.New(cfg, log, router)
	if err := server.Run(); err != nil {
		log.Error("server stopped with error", sl.Err(err))
	}
}
