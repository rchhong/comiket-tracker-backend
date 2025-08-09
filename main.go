package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/rchhong/comiket-backend/dao"
	"github.com/rchhong/comiket-backend/db"
)

func main() {
	mux := http.NewServeMux()

	postgresDB := db.InitializeDB()
	defer postgresDB.Teardown()

	userDAO := dao.NewUserDAO(postgresDB.Dbpool)

	mux.HandleFunc("GET /user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			log.Printf("Error converting path parameter %s: %s\n", id, err)
		}

		user, err := userDAO.GetUserByDiscordId(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusAccepted)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
		}

	})

	fmt.Printf("Listening on http://localhost:3000\n")
	log.Fatal(http.ListenAndServe(":3000", mux))

}
