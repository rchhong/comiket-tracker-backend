package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rchhong/comiket-backend/internal/config"
	"github.com/rchhong/comiket-backend/internal/controllers"
	"github.com/rchhong/comiket-backend/internal/db"
	"github.com/rchhong/comiket-backend/internal/logging"
	"github.com/rchhong/comiket-backend/internal/repositories/postgres"
	"github.com/rchhong/comiket-backend/internal/service"
	"github.com/rchhong/comiket-backend/internal/service/currency/ipgeoapi"
	"github.com/rchhong/comiket-backend/internal/service/scrape"
)

func main() {
	mux := http.NewServeMux()

	configPath := "/app/cmd/config.yaml"
	config, err := config.LoadConfigFromFile(configPath)
	if err != nil {
		log.Fatalf("Unable to load application configuration from file: %v", err)
	}

	err = logging.InitializeLogging(config.Logging.LogLevel, config.Logging.File.LogFilePath)
	if err != nil {
		log.Fatalf("Unable to setup logging: %v", err)
	} else {
		if config.Logging.File.LogFilePath == "" {
			logging.Logger.Warn("No logging file was supplied, logs will only be outputted to stdout")
		}
		logging.Logger.Info("Initialized logging successfully")
	}

	postgresDB, err := db.InitializeDB(config.Db.Postgres.Host, config.Db.Postgres.Port, config.Db.Postgres.DatabaseName, config.Db.Postgres.Username, config.Db.Postgres.Password)
	if err != nil {
		log.Fatalf("Unable to setup database: %v", err)
	} else {
		logging.Logger.Info("Initialized database successfully")
	}
	defer postgresDB.Teardown()

	currencyConverter, err := ipgeoapi.NewCurrencyConverterIpGeoAPI(ipgeoapi.IPGEO_API_CURRENCY_API_URL, os.Getenv("CURRENCY_API_KEY"), "JPY", "USD")
	if err != nil {
		log.Fatalf("Unable to retrieve currency conversion rate: %v", err)
	} else {
		logging.Logger.Info("Initialized currency converter successfully")
	}

	melonbooksScraper := scrape.NewMelonbooksScraper()

	userRepository := postgres.NewUserRepositoryPostgres(postgresDB)
	doujinRepository := postgres.NewDoujinRepositoryPostgres(postgresDB)
	reservationRepository := postgres.NewReservationRepositoryPostgres(postgresDB)
	exportRepository := postgres.NewExportRepositoryPostgres(postgresDB)

	currencyService := service.NewCurrencyService(currencyConverter)
	melonbooksScraperService := service.NewMelonbooksScraperService(melonbooksScraper, currencyService)
	userService := service.NewUserService(userRepository)
	doujinService := service.NewDoujinService(doujinRepository, melonbooksScraperService)
	reservationService := service.NewReservationService(reservationRepository, userService, doujinService)
	exportService := service.NewExportService(exportRepository)

	userController := controllers.NewUserController(userService)
	doujinController := controllers.NewDoujinController(doujinService)
	reservationController := controllers.NewReservationController(reservationService)
	exportController := controllers.NewExportController(exportService)

	userController.RegisterUserController(mux)
	doujinController.RegisterDoujinController(mux)
	reservationController.RegisterReservationController(mux)
	exportController.RegisterExportController(mux)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HEALTHY"))
	})

	logging.Logger.Info(fmt.Sprintf("Listening on http://localhost:%d", config.App.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.App.Port), mux))

}
