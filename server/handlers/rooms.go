package handlers

import (
	"net/http"
	"encoding/json"
	"log"
	"../models"
	"github.com/gorilla/mux"
)

func RoomViewSet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	id := vars["id"]
	defer r.Body.Close()

	switch r.Method  {
	case "GET":
		if id != "" {
			room := models.RetrieveRoomByID(id)
			err := json.NewEncoder(w).Encode(&room)
			if err != nil {
				log.Println("Unable to write JSON respone: ", err)
			}
			return
		}
		rooms := models.FetchAllRooms()
		err := json.NewEncoder(w).Encode(&rooms)
		if err != nil {
			log.Println("Unable to write JSON respone: ", err)
		}
	case "POST":
		var room models.Room
		err := json.NewDecoder(r.Body).Decode(&room)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Unable to read JSON request data: ", err)
			return
		}
		models.CreateRoom(&room)
		w.WriteHeader(http.StatusCreated)
	case "PUT":
		if id != "" {
			var room models.Room
			err := json.NewDecoder(r.Body).Decode(&room)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Println("Unable to read JSON request data: ", err)
				return
			}
			models.UpdateRoom(id, room)
		}
	case "DELETE":
		if id != "" {
			models.DeleteRoomByID(id)
		}
	}
}
