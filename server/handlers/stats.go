package handlers

import (
	"encoding/json"
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"../models"
)

func StatsViewSet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	id := vars["id"]
	defer r.Body.Close()

	switch r.Method {
	case "GET":
		if id != "" {
			statSetList := models.RetrieveGameStats(id)
			err := json.NewEncoder(w).Encode(&statSetList)
			if err != nil {
				log.Println("Unable to write JSON respone: ", err)
			}
			return
		}
		http.Error(w, "You can just retrieve details.", 400)
		return
	default:
		http.Error(w, "Invalid request method.", 405)
	}
}
