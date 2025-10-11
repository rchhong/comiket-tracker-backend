package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rchhong/comiket-backend/controllers"
	"github.com/rchhong/comiket-backend/dao"
	"github.com/rchhong/comiket-backend/db"
	"github.com/rchhong/comiket-backend/scrape"
	"github.com/rchhong/comiket-backend/service"
	"github.com/rchhong/comiket-backend/utils"
)

func main() {
	mux := http.NewServeMux()

	postgresDB := db.InitializeDB()
	defer postgresDB.Teardown()

	currencyConverter, err := scrape.NewCurrencyConverterImpl(utils.CURRENCY_API_URL, os.Getenv("CURRENCY_API_KEY"), "JPY", "USD")
	if err != nil {
		log.Fatalf("[ERROR] unable to retrieve currency conversion rate: %v", err)
	}

	melonbooksScraper := scrape.NewMelonbooksScraper(currencyConverter)
	melonbooksScraperService := service.NewMelonbooksScraperService(melonbooksScraper)

	userDao := dao.NewUserDAO(postgresDB.Dbpool)
	doujinDao := dao.NewDoujinDao(postgresDB.Dbpool)
	reservationDao := dao.NewReservationDAO(postgresDB.Dbpool)

	userService := service.NewUserService(userDao)
	doujinService := service.NewDoujinService(doujinDao, melonbooksScraperService)
	reservationService := service.NewReservationService(reservationDao, userService, doujinService)

	userController := controllers.NewUserController(userService, reservationService)
	doujinController := controllers.NewDoujinController(doujinService, reservationService)

	userController.RegisterUserController(mux)
	doujinController.RegisterDoujinController(mux)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HEALTHY"))
	})

	fmt.Printf("Listening on http://localhost:3000\n")
	log.Fatal(http.ListenAndServe(":3000", mux))

}
