package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rchhong/comiket-backend/dao"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	mux := http.NewServeMux()

	dao := dao.NewDAO(os.Getenv("DATABASE_URL"))
	defer dao.Close()

	mux.HandleFunc("GET /user/{id}", func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")

		id, err := primitive.ObjectIDFromHex(idString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to decode Id %s", idString), http.StatusBadRequest)
		}

		user, err := dao.RetrieveUserById(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to decode retrieve user with Id %s: %s", idString, err), http.StatusBadRequest)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Printf("user: %s", user)
		json.NewEncoder(w).Encode(user)

	})

	fmt.Printf("Listening on http://localhost:3000\n")
	log.Fatal(http.ListenAndServe(":3000", mux))

}
