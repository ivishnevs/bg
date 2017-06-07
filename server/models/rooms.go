package models

import (
	"time"
)

type Room struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Name        string `json:"name"`
	Description string `json:"description"`

	UserID uint
	Games  []Game `json:"games"`
}

type ExportedRoom struct {
	Room
	User User
}

func FetchAllRooms() []ExportedRoom {
	var rooms []Room
	var exported []ExportedRoom
	db.Find(&rooms)
	for _, room := range rooms {
		user := User{}
		db.Model(&room).Related(&user)
		exported = append(exported, ExportedRoom{Room: room, User:user})
	}
	return exported
}

func RetrieveRoomByID(id string) Room {
	room := Room{}
	games := []Game{}
	db.First(&room, id)
	db.Model(&room).Related(&games)

	for index, game := range games {
		gamers := []Gamer{}
		db.Model(&game).Related(&gamers)
		games[index].Gamers = gamers
	}


	room.Games = games
	return room
}

func CreateRoom(room *Room) {
	db.Create(room)
}

func UpdateRoom(id string, newRoom Room) {
	room := Room{}
	//db.First(&room, id)
	db.Model(&room).Where("id = ?", id).Updates(newRoom)
}

func DeleteRoomByID(id string) {
	if id != "" {
		db.Where("id = ?", id).Delete(&Room{})
	}
}
