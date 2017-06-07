package handlers

import (
	"net/http"
	"fmt"
	"sync"
	"crypto/rand"
	"crypto/md5"
	"../models"
	"encoding/json"
	"log"
	"time"
)

type SessionData struct {
	UserID uint
	ExpIn time.Time
}

type Session struct {
	sync.RWMutex
	Data map[string]SessionData
	Key string
	MaxAge int
}

func (s Session) generateKey() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (s Session) generateUniqueKey() string {
	key := s.generateKey()
	_, ok := s.Data[key]

	if ok {
		return s.generateUniqueKey()
	}
	return key
}

func (s Session) Set(data SessionData) string {
	s.Lock()
	defer s.Unlock()
	key := s.generateUniqueKey()
	s.Data[key] = data
	return key
}

func (s Session) Get(key string) (SessionData, bool) {
	s.RLock()
	defer s.RUnlock()
	data, ok := s.Data[key]
	return data, ok
}

func (s Session) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.Data, key)
}


var userSession Session

func init() {
	userSession = Session{
		Key: "sessionKey",
		MaxAge: 60 * 60 * 24,
		Data: make(map[string]SessionData),
	}
	go func() {
		for ; ; {
			time.Sleep(time.Duration(userSession.MaxAge/2) * time.Second)
			log.Println("Start cleaning userSession")
			for sessionKey, sessionData := range userSession.Data {
				if time.Since(sessionData.ExpIn) > 0 {
					log.Printf("Session with key %v\n for user id='%v' is expired\n", sessionKey, sessionData.UserID)
					userSession.Delete(sessionKey)
				}
			}
		}
	}()
}

func CurrentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	defer r.Body.Close()
	cookies := r.Cookies()

	var dataKey string
	user := models.User{}

	for _, cookie := range cookies {
		if cookie.Name == userSession.Key {
			dataKey = cookie.Value
			break
		}
	}
	sessionData, ok := userSession.Get(dataKey)
	if !ok {
		if err := json.NewEncoder(w).Encode(&user); err != nil {
			log.Println("Unable to write JSON respone: ", err)
		}
		return
	}
	user = models.GetUserByID(sessionData.UserID)

	if err := json.NewEncoder(w).Encode(&user); err != nil {
		log.Println("Unable to write JSON respone: ", err)
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := models.GetUserByEmail(r.PostFormValue("email"))
	if user.ID == 0 {
		http.Error(w, "No such user.", 404)
		return
	}
	passHash := fmt.Sprintf("%x", md5.Sum([]byte(r.PostFormValue("password"))))

	if user.PassHash != passHash {
		http.Error(w, "Incorrect password or email.", 400)
		return
	}

	dataKey := userSession.Set(SessionData{
		UserID: user.ID,
		ExpIn: time.Now().Local().Add(time.Duration(userSession.MaxAge)*time.Second),
	})
	cookie := http.Cookie{
		Name:     userSession.Key,
		Value:    dataKey,
		MaxAge:   userSession.MaxAge,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case "POST":
		user := models.User{
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			PassHash: fmt.Sprintf("%x", md5.Sum([]byte(r.PostFormValue("password")))),
		}

		models.CreateUser(&user)
		if user.ID != 0 {
			user.Room = models.Room{Name: r.PostFormValue("room")}
			user.Save()

		}

		dataKey := userSession.Set(SessionData{
			UserID: user.ID,
			ExpIn: time.Now().Local().Add(time.Duration(userSession.MaxAge)*time.Second),
		})
		cookie := http.Cookie{
			Name:     userSession.Key,
			Value:    dataKey,
			MaxAge:   userSession.MaxAge,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Invalid request method.", 405)
	}
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()

	var dataKey string
	for _, cookie := range cookies {
		if cookie.Name == userSession.Key {
			dataKey = cookie.Value
			break
		}
	}
	userSession.Delete(dataKey)

	cookie := http.Cookie{
		Name:     userSession.Key,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}
