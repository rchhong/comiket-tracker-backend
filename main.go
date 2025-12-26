package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rchhong/comiket-backend/controllers"
	"github.com/rchhong/comiket-backend/currency/ipgeoapi"
	"github.com/rchhong/comiket-backend/db"
	"github.com/rchhong/comiket-backend/repositories/postgres"
	"github.com/rchhong/comiket-backend/scrape"
	"github.com/rchhong/comiket-backend/service"
)

func main() {
	mux := http.NewServeMux()

	postgresDB := db.InitializeDB()
	defer postgresDB.Teardown()

	currencyConverter, err := ipgeoapi.NewCurrencyConverterIpGeoAPI(ipgeoapi.IPGEO_API_CURRENCY_API_URL, os.Getenv("CURRENCY_API_KEY"), "JPY", "USD")
	if err != nil {
		log.Fatalf("[ERROR] unable to retrieve currency conversion rate: %v", err)
	}

	melonbooksScraper := scrape.NewMelonbooksScraper(currencyConverter)

	userRepository := postgres.NewUserRepositoryPostgres(postgresDB.Dbpool)
	doujinRepository := postgres.NewDoujinRepositoryPostgres(postgresDB.Dbpool)
	reservationRepository := postgres.NewReservationRepositoryPostgres(postgresDB.Dbpool)
	exportRepository := postgres.NewExportRepositoryPostgres(postgresDB.Dbpool)

	userService := service.NewUserService(userRepository)
	doujinService := service.NewDoujinService(doujinRepository, melonbooksScraper)
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
