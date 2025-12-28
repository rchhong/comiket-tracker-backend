package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rchhong/comiket-backend/internal/controllers"
	"github.com/rchhong/comiket-backend/internal/db"
	"github.com/rchhong/comiket-backend/internal/repositories/postgres"
	"github.com/rchhong/comiket-backend/internal/service"
	"github.com/rchhong/comiket-backend/internal/service/currency/ipgeoapi"
	"github.com/rchhong/comiket-backend/internal/service/scrape"
)

func main() {
	mux := http.NewServeMux()

	postgresDB, err := db.InitializeDB()
	if err != nil {
		log.Fatalf("Unable to setup database: %s", err)
	}
	defer postgresDB.Teardown()

	currencyConverter, err := ipgeoapi.NewCurrencyConverterIpGeoAPI(ipgeoapi.IPGEO_API_CURRENCY_API_URL, os.Getenv("CURRENCY_API_KEY"), "JPY", "USD")
	if err != nil {
		log.Fatalf("[ERROR] unable to retrieve currency conversion rate: %v", err)
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

	fmt.Printf("Listening on http://localhost:3000\n")
	log.Fatal(http.ListenAndServe(":3000", mux))

}
