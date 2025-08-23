package main

import (
	"time"
	"math/rand"
	"context"
	"sync"
)

type Person struct {
	Id       int
	Name     string
	BirthDay time.Time
	GenderID   int
	MainRaceID int
	SubRaceIDs []int
	MainJobID  int
	SubJobIDs  []int
	Level    int
	Alive    bool
	DeathDay time.Time
}


func NewPerson(name string) Person {
	LastId.Person++
	var genderID int
	var raceID int
	var jobID int
	if len(World.Genders) > 0 {
		genderID = World.Genders[rand.Intn(len(World.Genders))].Id
	}
	if len(World.Races) > 0 {
		raceID = World.Races[rand.Intn(len(World.Races))].Id
	}
	if len(World.Jobs) > 0 {
		jobID = World.Jobs[rand.Intn(len(World.Jobs))].Id
	}
	return Person{
		Id:       LastId.Person,
		Name:     name,
		BirthDay: time.Now(),
		GenderID:   genderID,
		MainRaceID: raceID,
		SubRaceIDs: []int{},
		MainJobID:  jobID,
		SubJobIDs:  []int{},
		Level:    1,
		Alive:    true,
		DeathDay: time.Time{},
	}
}

// Run は住人の行動ループ。状態は直接変更せず、イベントを送信する。
func (p Person) Run(ctx context.Context, areaID int, events chan<- Event, wg *sync.WaitGroup) {
    defer wg.Done()
    // 0.5s〜1.5sの間隔で行動
    jitter := time.Duration(250+rand.Intn(500)) * time.Millisecond
    ticker := time.NewTicker(750*time.Millisecond + jitter)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            // 行動終了イベント（レベルアップと収入加算は適用ループで実行）
            events <- PersonActed{PersonID: p.Id, AreaID: areaID}
        }
    }
}
