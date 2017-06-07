package models

import (
	"time"
	"encoding/json"
	"log"
	"fmt"
)

type Gamer struct {
	ID         uint `json:"id" gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
	LastAction time.Time `json:"lastAction"`

	Role            int `json:"role"`
	Supplier        *Gamer `gorm:"ForeignKey:SupplierID;AssociationForeignKey:Supplier" json:"-"`
	SupplierID      uint `json:"-"`
	Customer        *Gamer `gorm:"ForeignKey:CustomerID;AssociationForeignKey:Customer" json:"-"`
	CustomerID      uint `json:"-"`
	OrderQueue      string `gorm:"default:'[6]'" json:"-"`
	SupplyQueue     string `gorm:"default:'[6]'" json:"-"`
	Storage         int `gorm:"default:'10'" json:"storage"`
	CurrentOrder    int `gorm:"default:'6'" json:"currentOrder"`
	GamerOrder      int `gorm:"default:'6'" json:"gamerOrder"`
	Debt            int `gorm:"default:'0'" json:"debt"`
	Penalty         float32 `json:"penalty"`
	IsActive        bool `gorm:"default:'false'" json:"isActive"`
	IsStepCompleted bool `gorm:"default:'false'" json:"isStepCompleted"`

	Stats  []Stats
	GameID uint `gorm:"index"`
}

func JsonDumps(a interface{}) string {
	b, err := json.Marshal(a)
	if err != nil {
		log.Println("Unable to Marshal JSON: ", err)
	}
	return string(b)
}

func JsonLoads(s string, a interface{}) {
	err := json.Unmarshal([]byte(s), a)
	if err != nil {
		log.Println("Unable to Unmarshal JSON: ", err)
	}
}

func RetrievGamerByID(id string) *Gamer {
	gamer := Gamer{}
	db.First(&gamer, id)
	return &gamer
}

type GamerMetadata struct {
	GameID      uint `json:"gameID"`
	CurrentStep int `json:"currentStep"`
	StepsNumber int `json:"stepsNumber"`
	GamerCount  int `json:"gamerCount"`
	Gamers      []Gamer `json:"gamers"`
}

type GameData struct {
	GamerMetadata `json:"gamerMetadata"`
	Stats    []Stats `json:"stats"`
	RoomName string `json:"roomName"`
	Gamer
}

func RetrieveGameplayFlow(id string) (interface{}, bool) {
	gamer := Gamer{}
	room := Room{}
	db.First(&gamer, id)

	game := Game{}
	db.Model(&gamer).Related(&game)

	db.Model(&game).Related(&room)

	stats := []Stats{}
	db.Model(&gamer).Related(&stats)

	if game.StepsNumber < game.CurrentStep { // game is finished
		return game.ID, true
	}

	gamers := []Gamer{}
	db.Model(&game).Related(&gamers)

	metadata := GamerMetadata{
		GameID:      game.ID,
		CurrentStep: game.CurrentStep,
		StepsNumber: game.StepsNumber,
		GamerCount:  game.GamerCount,
		Gamers:      gamers,
	}
	data := GameData{
		GamerMetadata: metadata,
		Stats:         stats,
		RoomName:      room.Name,
		Gamer:         gamer,
	}
	return data, false
}

func ActivateRole(id string) bool {
	gamer := Gamer{}
	db.First(&gamer, id)
	if gamer.IsActive {
		return false
	}
	db.Model(&gamer).Updates(map[string]interface{}{
		"IsActive": true,
		"LastAction": time.Now(),
	})
	game := Game{}
	db.Model(&gamer).Related(&game)

	newOccupiedPlaces := game.OccupiedPlaces + 1
	if newOccupiedPlaces >= game.GamerCount {
		newOccupiedPlaces = game.GamerCount
	} else if newOccupiedPlaces <= 0 {
		newOccupiedPlaces = 0
	}

	db.Model(&game).Updates(map[string]interface{}{
		"OccupiedPlaces": newOccupiedPlaces,
	})
	return true
}

func ReleaseNonActiveGamers(timeout time.Duration)  {
	gamers := []Gamer{}
	db.Find(&gamers)
	fmt.Println("ReleaseNonActiveGamers")
	for _, gamer := range gamers {
		var gamerActiveTimeout time.Duration = timeout
		if gamer.IsStepCompleted {
			gamerActiveTimeout = timeout * time.Duration(2)
		}
		fmt.Println(gamerActiveTimeout)
		if time.Since(gamer.LastAction) > gamerActiveTimeout && gamer.IsActive {
			fmt.Println(gamer.Role)
			game := Game{}
			db.Model(&gamer).Updates(map[string]interface{}{
				"IsActive": false,
			})
			db.Model(&gamer).Related(&game)
			newOccupiedPlaces := game.OccupiedPlaces - 1
			if newOccupiedPlaces >= game.GamerCount {
				newOccupiedPlaces = game.GamerCount
			} else if newOccupiedPlaces <= 0 {
				newOccupiedPlaces = 0
			}
			db.Model(&game).Updates(map[string]interface{}{
				"OccupiedPlaces": newOccupiedPlaces,
			})
		}
	}
}

func (gamer *Gamer) PushToOrderQueue(value int) {
	orderQueue := []int{}
	JsonLoads(gamer.OrderQueue, &orderQueue)
	orderQueue = append(orderQueue, value)
	gamer.OrderQueue = JsonDumps(orderQueue)
	db.Model(&gamer).Updates(gamer)

}

func (gamer *Gamer) PushToSupplyQueue(value int) {
	supplyQueue := []int{}
	JsonLoads(gamer.SupplyQueue, &supplyQueue)
	supplyQueue = append(supplyQueue, value)
	gamer.SupplyQueue = JsonDumps(supplyQueue)
	db.Model(&gamer).Updates(gamer)
}

func (gamer *Gamer) PopFromOrderQueue() int {
	orderQueue := []int{}
	JsonLoads(gamer.OrderQueue, &orderQueue)
	firstElement := orderQueue[0]
	orderQueue = orderQueue[1:]
	gamer.OrderQueue = JsonDumps(orderQueue)
	db.Model(&gamer).Updates(gamer)
	return firstElement
}

func (gamer *Gamer) PopFromSupplyQueue() int {
	supplyQueue := []int{}
	JsonLoads(gamer.SupplyQueue, &supplyQueue)
	firstElement := supplyQueue[0]
	supplyQueue = supplyQueue[1:]
	gamer.SupplyQueue = JsonDumps(supplyQueue)
	db.Model(&gamer).Updates(gamer)
	return firstElement
}

func (gamer *Gamer) MakeOrder(value int) {
	if gamer.SupplierID != 0 {
		supplier := Gamer{}
		db.First(&supplier, gamer.SupplierID)
		supplier.PushToOrderQueue(value)
	} else {
		gamer.PushToSupplyQueue(value)
	}
	if gamer.CustomerID == 0 {
		game := Game{}
		db.Model(&gamer).Related(&game)
		demand := game.DemandFunc()
		gamer.PushToOrderQueue(demand)
	}

	db.Model(&gamer).Updates(map[string]interface{}{
		"GamerOrder": value,
	})
}

func (gamer *Gamer) MakeSupply() {
	game := Game{}
	db.Model(&gamer).Related(&game)

	var supply int
	order := gamer.CurrentOrder
	debt := gamer.Debt

	if gamer.Storage >= order + debt {
		supply = order + debt
		gamer.Penalty += float32(gamer.Storage - supply) * game.HoldingCost
		db.Create(&Stats{
			Storage:      gamer.Storage,
			CurrentOrder: gamer.CurrentOrder,
			Debt:	      gamer.Debt,
			GamerOrder:   gamer.GamerOrder,
			Penalty:      gamer.Penalty,
			Step:         game.CurrentStep,
			GamerID:      gamer.ID,
		})

		gamer.Storage -= supply
		gamer.CurrentOrder = 0
		gamer.Debt = 0
	} else {
		supply = gamer.Storage
		gamer.Penalty += float32(order + debt - supply) * game.BackorderCost
		db.Create(&Stats{
			Storage:      gamer.Storage,
			CurrentOrder: gamer.CurrentOrder,
			Debt:	      gamer.Debt,
			GamerOrder:   order,
			Penalty:      gamer.Penalty,
			Step:         game.CurrentStep,
			GamerID:      gamer.ID,
		})

		gamer.Debt = order + debt - gamer.Storage
		gamer.CurrentOrder = 0
		gamer.Storage = 0
	}
	if gamer.CustomerID != 0 {
		customer := Gamer{}
		db.First(&customer, gamer.CustomerID)
		customer.PushToSupplyQueue(supply)
	}
	db.Model(&gamer).Updates(map[string]interface{}{
		"Storage": gamer.Storage,
		"CurrentOrder": gamer.CurrentOrder,
		"Debt": gamer.Debt,
		"Penalty": gamer.Penalty,
	})
}

func (gamer *Gamer) CompleteStep() {
	gamer.IsStepCompleted = true
	gamer.LastAction = time.Now()
	db.Model(&gamer).Updates(gamer)

	game := &Game{}
	db.Model(&gamer).Related(&game)
	game.ProcessStepCompletion()
}

func (gamer *Gamer) DumpsStats(order int) {
	game := &Game{}
	db.Model(&gamer).Related(&game)
	db.Create(&Stats{
		Storage:      gamer.Storage,
		CurrentOrder: gamer.CurrentOrder,
		Debt:	      gamer.Debt,
		GamerOrder:   order,
		Penalty:      gamer.Penalty,
		Step:         game.CurrentStep,
		GamerID:      gamer.ID,
	})
}

func (gamer *Gamer) PerformGameFlow(order int) {
	gamer.MakeOrder(order)
	gamer.CompleteStep()
}
