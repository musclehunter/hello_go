package main

import (
	"fmt"
	"os"
	"time"
	"net/http"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	gofakeit.Seed(time.Now().UnixNano())
	config, err := LoadConfig("config.yml")
	if err != nil {
		fmt.Println("config.yml not found", err)
		os.Exit(1)
	}

	world := loadWorld()
	if world.Id == 0 {
		world = NewWorld("first world")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(world)
	})

	go func() {
		addr := fmt.Sprintf(":%d", config.Port)
		fmt.Printf("listening on %s\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Println("ListenAndServe: ", err)
		}
	}()

	count := 1
	last := time.Now()
	for {
		now := time.Now()
		if now.Sub(last) >= time.Duration(config.TurnSeconds) * time.Second {
			action(&world)
			last = now
			count++
		}
		saveWorld(world)
		saveLastId(lastId)
		time.Sleep(100 * time.Millisecond)
	}
}

func action(world *World) {
	// areaが空だったら作る
	if len(world.Areas) == 0 {
		world.Areas = append(world.Areas, NewArea(gofakeit.BeerName()))
	}
	// personが空だったら作る
	if len(world.Persons) == 0 {
		world.Persons = append(world.Persons, NewPerson(gofakeit.Name()))
	}
	// genderが空だったら作る
	if len(world.Genders) == 0 {
		world.Genders = append(world.Genders, NewGender(gofakeit.Gender()))
	}
	// raceが空だったら作る
	if len(world.Races) == 0 {
		world.Races = append(world.Races, NewRace(gofakeit.Color() + gofakeit.Animal()))
	}
	// jobが空だったら作る
	if len(world.Jobs) == 0 {
		world.Jobs = append(world.Jobs, NewJob(gofakeit.JobTitle() + gofakeit.Dessert()))
	}	
	
	jsonBytes, err := json.MarshalIndent(world, "", " ")
	if err != nil {
		fmt.Println("json.MarshalIndent: ", err)
		return
	} else {
		fmt.Printf("[%s]\n%s\n", time.Now().Format("2006-01-02 15:04:05"), string(jsonBytes))
	}
	os.Stdout.Sync()
}