package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rchhong/comiket-backend/controllers"
	"github.com/rchhong/comiket-backend/dao"
	"github.com/rchhong/comiket-backend/db"
	"github.com/rchhong/comiket-backend/service"
)

func main() {
	mux := http.NewServeMux()

	postgresDB := db.InitializeDB()
	defer postgresDB.Teardown()

	userDAO := dao.NewUserDAO(postgresDB.Dbpool)
	userService := service.NewUserService(userDAO)
	userController := controllers.NewUserController(userService)

	userController.RegisterUserController(mux)

	fmt.Printf("Listening on http://localhost:3000\n")
	log.Fatal(http.ListenAndServe(":3000", mux))

}
