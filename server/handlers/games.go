package handlers

import (
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"../models"
)

func GameViewSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	defer r.Body.Close()

	switch r.Method  {
	case "GET":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if id != "" {
			game := models.RetrieveGameByID(id)
			if game.StepsNumber <= game.CurrentStep {
				redirectData := struct {
					IsGameFinished bool `json:"isGameFinished"`
					GameID interface{} `json:"gameId"`
				} {
					IsGameFinished: true,
					GameID: game.ID,
				}
				err := json.NewEncoder(w).Encode(&redirectData)
				if err != nil {
					log.Println("Unable to write JSON respone: ", err)
				}
				return
			}
			if err := json.NewEncoder(w).Encode(&game); err != nil {
				log.Println("Unable to write JSON respone: ", err)
			}
			return
		}
		http.Error(w, "You can just retrieve details.", 400)
		return
	case "POST":
		game := models.Game{Status: "open"}
		if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Unable to read JSON request data: ", err)
			return
		}

		/////
		var dataKey string
		for _, cookie := range r.Cookies() {
			if cookie.Name == userSession.Key {
				dataKey = cookie.Value
				break
			}
		}
		sessionData, ok := userSession.Get(dataKey)
		if !ok {
			http.Error(w, "not authorized", 401)
			return
		}
		user := models.GetUserByID(sessionData.UserID)
		if user.ID == 0 {
			http.Error(w, "not authorized", 401)
			return
		}
		if int(user.Room.ID) != game.RoomID {
			http.Error(w, "Forbidden", 403)
			return
		}
		///////


		models.CreateGame(&game)
		game.CreateGamers()
		w.WriteHeader(http.StatusCreated)
	case "PUT":
		if id != "" {
			game := models.RetrieveGameByID(id)

			/////
			var dataKey string
			for _, cookie := range r.Cookies() {
				if cookie.Name == userSession.Key {
					dataKey = cookie.Value
					break
				}
			}
			sessionData, ok := userSession.Get(dataKey)
			if !ok {
				http.Error(w, "not authorized", 401)
				return
			}
			user := models.GetUserByID(sessionData.UserID)
			if user.ID == 0 {
				http.Error(w, "not authorized", 401)
				return
			}
			if int(user.Room.ID) != game.RoomID {
				http.Error(w, "Forbidden", 403)
				return
			}
			///////

			gameUpdatedFields := map[string]interface{}{}
			if err := json.NewDecoder(r.Body).Decode(&gameUpdatedFields); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Println("Unable to read JSON request data: ", err)
				return
			}
			_, isReset := gameUpdatedFields["occupiedPlaces"]
			isInGame := game.CurrentStep > 1
			needGamersReproduce := true

			if !isReset && isInGame {
				needGamersReproduce = false
				delete(gameUpdatedFields, "gamerCount")
			}

			models.UpdateGame(id, gameUpdatedFields)
			if needGamersReproduce {
				game := models.RetrieveGameByID(id)
				game.CreateGamers()
			}
			w.WriteHeader(http.StatusOK)
		}
	case "DELETE":
		if id != "" {
			game := models.RetrieveGameByID(id)

			/////
			var dataKey string
			for _, cookie := range r.Cookies() {
				if cookie.Name == userSession.Key {
					dataKey = cookie.Value
					break
				}
			}
			sessionData, ok := userSession.Get(dataKey)
			if !ok {
				http.Error(w, "not authorized", 401)
				return
			}
			user := models.GetUserByID(sessionData.UserID)
			if user.ID == 0 {
				http.Error(w, "not authorized", 401)
				return
			}
			if int(user.Room.ID) != game.RoomID {
				http.Error(w, "Forbidden", 403)
				return
			}
			///////

			models.DeleteGameByID(id)
		}
	}
}
