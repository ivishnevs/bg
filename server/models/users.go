package models

import (
	"time"
)

type User struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Name     string `json:"name"`
	Email    string `gorm:"not null;unique" json:"email"`
	PassHash string `gorm:"not null;" json:"-"`

	Room Room
}

func (u *User) Save()  {
	db.Save(u)
}

func CreateUser(user *User) {
	db.Create(user)
}

func GetUserByID(id uint) User {
	user := User{}
	room := Room{}
	db.First(&user, id)

	db.Model(&user).Related(&room)
	user.Room = room
	return user
}

func GetUserByEmail(email string) User {
	user := User{}
	room := Room{}
	db.Where("email = ?", email).First(&user)

	db.Model(&user).Related(&room)
	user.Room = room
	return user
}
