package handlers

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"sync"
)

type subscriber *websocket.Conn
type subscribersSet map[subscriber]bool

type Action struct {
	ActionType string `json:"type"`
	Data       string `json:"data"`
}

type Event struct {
	Channel string `json:"channel"`
	Action  Action `json:"action"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var subscribersByChannel = struct {
	sync.RWMutex
	m map[string]subscribersSet
}{m: make(map[string]subscribersSet)}

func RegisterSub(sub subscriber, channel string)  {
	subscribersByChannel.Lock()
	defer subscribersByChannel.Unlock()
	log.Printf("Registering subscriber to channel %v\n", channel)

	if _, ok := subscribersByChannel.m[channel]; !ok {
		log.Printf("Creating channel %v\n", channel)
		subscribersByChannel.m[channel] = make(subscribersSet)
	}
	subscribersByChannel.m[channel][sub] = true
}

func UnregisterSub(sub subscriber, channel string) {
	subscribersByChannel.Lock()
	defer subscribersByChannel.Unlock()
	log.Printf("Unregistering subscriber from channel %v\n", channel)
	subs, ok := subscribersByChannel.m[channel]
	if !ok {
		log.Printf("There are no '%v' channel \n", channel)
		return
	}
	if _, ok = subs[sub]; ok {
		delete(subs, sub)
	}
}

func RemoveSub(sub subscriber) {
	subscribersByChannel.Lock()
	defer subscribersByChannel.Unlock()
	log.Println("Removing subscriber")
	for _, subs := range subscribersByChannel.m {
		if _, ok := subs[sub]; ok {
			delete(subs, sub)
		}
	}
}

func NotifySubscribers(event Event) {
	subscribersByChannel.RLock()
	defer subscribersByChannel.RUnlock()
	log.Printf("Notifying subscribers about %v", event)
	for sub := range subscribersByChannel.m[event.Channel] {
		if sub != nil {
			if err := (*websocket.Conn)(sub).WriteJSON(event); err != nil {
				log.Fatalln("Unable to write JSON respone: ", err)
				break // жесткий выход
			}
		}
	}
}

func processWSEvent(sub subscriber, event Event)  {
	log.Printf("Receive event:\nChannel: %v\nAction: %v\nData: %v\n",
		event.Channel, event.Action.ActionType, event.Action.Data)
	switch event.Action.ActionType {
	case "subscribe.to":
		RegisterSub(sub, event.Action.Data)
		(*websocket.Conn)(sub).WriteJSON(event)
	case "role.selected":
		NotifySubscribers(event)
	case "unsubscribe.from":
		UnregisterSub(sub, event.Action.Data)
	}
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("open ws connection")
	defer log.Println("close ws connection")

	conn, err := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	defer RemoveSub(conn)

	if err != nil {
		log.Fatalln(err)
		return
	}

	for {
		event := Event{}
		if err := conn.ReadJSON(&event); err != nil {
			log.Println("Unable to read JSON request: ", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Bad request"))
			return
		}
		processWSEvent(conn, event)
	}
}
