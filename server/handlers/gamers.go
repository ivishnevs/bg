package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"log"
	"../models"
	"fmt"
	"strconv"
	"crypto/rand"
	"time"
	"sync"
)

var makeOrderMutex = &sync.Mutex{}

type orderMsg struct {
	Order       int `json:"order"`
}

type GamerSessionData struct {
	GamerID string
	ExpIn time.Time
}

type GamerSession struct {
	sync.RWMutex
	Data map[string]GamerSessionData
	Key string
	MaxAge int
}

func (s GamerSession) generateKey() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (s GamerSession) generateUniqueKey() string {
	key := s.generateKey()
	_, ok := s.Data[key]

	if ok {
		return s.generateUniqueKey()
	}
	return key
}

func (s GamerSession) Set(data GamerSessionData) string {
	s.Lock()
	defer s.Unlock()
	key := s.generateUniqueKey()
	s.Data[key] = data
	return key
}

func (s GamerSession) Get(key string) (GamerSessionData, bool) {
	s.RLock()
	defer s.RUnlock()
	data, ok := s.Data[key]
	return data, ok
}

func (s GamerSession) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.Data, key)
}


var gamerSession GamerSession

func init()  {
	gamerSession = GamerSession{
		Key: "gamerKey",
		MaxAge: 60 * 10,
		Data: make(map[string]GamerSessionData),
	}
	go func() {
		for ; ; {
			time.Sleep(3 * 60 * time.Second)
			models.ReleaseNonActiveGamers(time.Duration(gamerSession.MaxAge) * time.Second)
			for sessionKey, sessionData := range gamerSession.Data {
				if time.Since(sessionData.ExpIn) > 0 {
					log.Printf("Session with key %v\n for gamer id='%v' is expired\n", sessionKey, sessionData.GamerID)
					gamerSession.Delete(sessionKey)
				}
			}
		}
	}()
}

func setGamerSession(w http.ResponseWriter, gamerID string) {
	dataKey := gamerSession.Set(GamerSessionData{
		GamerID: gamerID,
		ExpIn: time.Now().Local().Add(time.Duration(gamerSession.MaxAge)*time.Second),
	})
	cookie := http.Cookie{
		Name:     gamerSession.Key,
		Value:    dataKey,
		MaxAge:   gamerSession.MaxAge,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	gameIdCookie := http.Cookie{
		Name:     "gamerID",
		Value:    gamerID,
		MaxAge:   gamerSession.MaxAge,
		Path:     "/",
		HttpOnly: false,
	}
	http.SetCookie(w, &gameIdCookie)
}

func checkAndDeleteGamerSession(r *http.Request, requestedID string, needClean bool) bool {
	cookies := r.Cookies()
	var dataKey string
	for _, cookie := range cookies {
		if cookie.Name == gamerSession.Key {
			dataKey = cookie.Value
			break
		}
	}
	gamerSessionData, ok := gamerSession.Get(dataKey)
	if !ok {
		log.Println("No gamerSessionData !! ")
		log.Println(dataKey)
		return false
	}
	if gamerSessionData.GamerID != requestedID {
		log.Println("gamerSessionData !! ", gamerSessionData.GamerID, requestedID)
		return false
	}
	if needClean {
		gamerSession.Delete(dataKey)
	}
	return true
}


func GamerViewSet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	id := vars["id"]
	defer r.Body.Close()

	if id == "" {
		http.Error(w, "You can use just details endpoind.", 400)
		return
	}

	switch r.Method {
	case "GET":
		isValidGamer := true
		if ok := models.ActivateRole(id); ok {
			setGamerSession(w, id)
		} else {
			isValidGamer = checkAndDeleteGamerSession(r, id, false)
		}
		if !isValidGamer {
			err := json.NewEncoder(w).Encode(&map[string]string{"redirect_url": "#/rooms"})
			if err != nil {
				log.Println("Unable to write JSON respone: ", err)
			}
			return
		}

		data, isGameFinished := models.RetrieveGameplayFlow(id)
		if isGameFinished {
			redirectData := struct {
				IsGameFinished bool `json:"isGameFinished"`
				GameID interface{} `json:"gameId"`
			} {
				IsGameFinished: true,
				GameID: data, // data here is gameID
			}
			err := json.NewEncoder(w).Encode(&redirectData)
			if err != nil {
				log.Println("Unable to write JSON respone: ", err)
			}
		} else {
			err := json.NewEncoder(w).Encode(&data)
			if err != nil {
				log.Println("Unable to write JSON respone: ", err)
			}
		}
	case "POST":
		if ok := checkAndDeleteGamerSession(r, id, false); !ok {
			err := json.NewEncoder(w).Encode(&map[string]string{"redirect_url": "#/rooms"})
			if err != nil {
				log.Println("Unable to write JSON respone: ", err)
			}
			return
		}
		//setGamerSession(w, id)
		msg := orderMsg{}
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Unable to read JSON request data: ", err)
			return
		}

		makeOrderMutex.Lock()
		gamer := models.RetrievGamerByID(id)
		if gamer.IsStepCompleted {
			http.Error(w, "Current step is already completed by this gamer.", 400)
			return
		}
		if msg.Order < 0 {
			http.Error(w, "Invalid order. Order must be a positive integer.", 400)
			return
		}
		gamer.PerformGameFlow(msg.Order)
		log.Println("Realize the lock")
		makeOrderMutex.Unlock()

		// notification
		channel := fmt.Sprintf("game-%v", gamer.GameID)
		event := Event{
			Channel: channel,
			Action: Action{
				ActionType: "gamer.step-completed",
				Data: strconv.Itoa(gamer.Role),
			},
		}
		NotifySubscribers(event)
		checkAndDeleteGamerSession(r, id, true)
		setGamerSession(w, id)
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Invalid request method.", 405)
	}
}
