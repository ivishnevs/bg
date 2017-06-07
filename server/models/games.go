package models

import (
	"time"
)

type Game struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	RoomID int `gorm:"index"`

	Status         string `json:"status"`
	GamerCount     int `json:"gamerCount"`
	OccupiedPlaces int `json:"occupiedPlaces"`
	StepsNumber    int `json:"stepsNumber"`
	CurrentStep    int `gorm:"default:'1'" json:"currentStep"`
	TotalPenalty   float32 `json:"totalPenalty"`
	HoldingCost    float32 `json:"holdingCost"`
	BackorderCost  float32 `json:"backorderCost"`
	DemandPattern  int `json:"demandPattern"`

	Gamers []Gamer `json:"gamers"`
}

func (game *Game) ProcessStepCompletion()  {
	isGameStepCompleted := true
	gamers := []Gamer{}
	db.Model(&game).Related(&gamers)

	for _, gamer := range gamers {
		if !gamer.IsStepCompleted {
			isGameStepCompleted = false
		}
	}
	if isGameStepCompleted {
		var TotalPenalty float32

		for _, gamer := range gamers {
			gamer.MakeSupply()
		}

		db.Model(&game).Related(&gamers)
		for _, gamer := range gamers {
			gamer.CurrentOrder += gamer.PopFromOrderQueue()
			gamer.Storage += gamer.PopFromSupplyQueue()
			gamer.IsStepCompleted = false
			TotalPenalty += gamer.Penalty

			db.Model(&gamer).Updates(map[string]interface{}{
				"Storage": gamer.Storage,
				"CurrentOrder": gamer.CurrentOrder,
				"IsStepCompleted": gamer.IsStepCompleted,
			})
		}
		game.TotalPenalty = TotalPenalty
		game.CurrentStep++
		db.Model(&game).Updates(game)
	}
}

func (game Game) CreateGamers() {
	gamers := []Gamer{}
	db.Model(&game).Related(&gamers)

	for _, gamer := range gamers {
		db.Delete(&gamer)
	}

	gamers = []Gamer{}

	for i := 0; i < game.GamerCount; i++ {
		gamer := Gamer{
			Role: i,
			GameID: game.ID,
			LastAction: time.Now(),
		}
		db.Create(&gamer)
		gamers = append(gamers, gamer)
	}

	for i := 1; i < game.GamerCount; i++ {
		prevGamer := gamers[i-1]
		gamer := gamers[i]
		prevGamer.SupplierID = gamer.ID
		gamer.CustomerID = prevGamer.ID
		db.Model(&gamer).Updates(gamer)
		db.Model(&prevGamer).Updates(prevGamer)
	}
}

func (game Game) DemandFunc() int {
	step := game.CurrentStep
	if game.DemandPattern == 1 {
		if step > 6 {
			return 8
		}
		if step < 4 {
			return 6
		}
		if step == 4 {
			return 8
		}
		if step == 5 {
			return 11
		}
		if step == 6 {
			return 12
		}
	}
	return 8
}

func RetrieveGameByID(id string) Game {
	game := Game{}
	gamers := []Gamer{}
	db.First(&game, id)

	db.Model(&game).Related(&gamers)
	game.Gamers = gamers
	return game
}

type StatsSet struct {
	GamerID uint `json:"gamerId"`
	GamerRole int `json:"gamerRole"`
	Stats []Stats `json:"stats"`
}

func RetrieveGameStats(id string) []StatsSet {
	game := Game{}
	gamers := []Gamer{}
	db.First(&game, id)
	db.Model(&game).Related(&gamers)

	statSetList := []StatsSet{}

	for _, gamer := range gamers {
		stats := []Stats{}
		db.Model(&gamer).Related(&stats)

		statSetList = append(statSetList, StatsSet{
			GamerID: gamer.ID,
			GamerRole: gamer.Role,
			Stats: stats,
		})

	}
	return statSetList
}

func CreateGame(game *Game)  {
	db.Create(game)
}

func UpdateGame(id string, gameUpdatedFields map[string]interface{})  {
	db.Model(&Game{}).Where("id = ?", id).Updates(gameUpdatedFields)
}

func DeleteGameByID(id string)  {
	if id != "" {
		game := Game{}
		db.Where("id = ?", id).First(&game)
		gamers := []Gamer{}
		db.Model(&game).Related(&gamers)

		for _, gamer := range gamers {
			db.Delete(&gamer)
		}
		if game.ID != 0 {
			db.Delete(&game)
		}
	}
	// TODO(ilya): check: delete all gamers
}
