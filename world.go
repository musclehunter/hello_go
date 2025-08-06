package main

import (
	"time"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"fmt"
)

type World struct {
	Id int
	Name string
	BirthDay time.Time
	Areas []Area
	Persons []Person
	Genders []Gender
	Races []Race
	Jobs []Job
}

type Area struct {
	Id int
	Name string
	Population int
	Residents []Person
}

type Person struct {
	Id int
	Name string
	BirthDay time.Time
	Gender Gender
	MainRace Race
	SubRace []Race
	MainJob Job
	SubJob []Job
}

type Gender struct {
	Id int
	Name string
}

type Race struct {
	Id int
	Name string
}

type Job struct {
	Id int
	Name string
}

type LastId struct {
	World int
	Area int
	Person int
	Gender int
	Race int
	Job int
}

var lastId LastId = loadLastId()

func NewWorld(name string) World {
	lastId.World++
	return World{
		Id: lastId.World,
		Name: name,
		BirthDay: time.Now(),
		Areas: []Area{},
		Persons: []Person{},
		Genders: []Gender{},
		Races: []Race{},
		Jobs: []Job{},
	}
}

func NewArea(name string) Area {
	lastId.Area++
	return Area{
		Id: lastId.Area,
		Name: name,
		Population: 0,
		Residents: []Person{},
	}
}

func NewPerson(name string) Person {
	lastId.Person++
	return Person{
		Id: lastId.Person,
		Name: name,
		BirthDay: time.Now(),
		Gender: Gender{},
		MainRace: Race{},
		SubRace: []Race{},
		MainJob: Job{},
		SubJob: []Job{},
	}
}

func NewGender(name string) Gender {
	lastId.Gender++
	return Gender{
		Id: lastId.Gender,
		Name: name,
	}
}

func NewRace(name string) Race {
	lastId.Race++
	return Race{
		Id: lastId.Race,
		Name: name,
	}
}

func NewJob(name string) Job {
	lastId.Job++
	return Job{
		Id: lastId.Job,
		Name: name,
	}
}

// Worldをymlに保存する
func saveWorld(world World) {
	data, err := yaml.Marshal(world)
	if err != nil {
		fmt.Println("yaml.Marshal: ", err)
		return
	}
	ioutil.WriteFile("world.yml", data, 0644)
}

func loadWorld() World {
	data, err := ioutil.ReadFile("world.yml")
	if err != nil {
		fmt.Println("ioutil.ReadFile: ", err)
		return World{}
	}
	var world World
	err = yaml.Unmarshal(data, &world)
	if err != nil {
		fmt.Println("yaml.Unmarshal: ", err)
		return World{}
	}
	return world
}

func saveLastId(lastId LastId) {
	data, err := yaml.Marshal(lastId)
	if err != nil {
		fmt.Println("yaml.Marshal: ", err)
		return
	}
	ioutil.WriteFile("last_id.yml", data, 0644)
}

func loadLastId() LastId {
	data, err := ioutil.ReadFile("last_id.yml")
	if err != nil {
		fmt.Println("ioutil.ReadFile: ", err)
		return LastId{
			World: 0,
			Area: 0,
			Person: 0,
			Gender: 0,
			Race: 0,
			Job: 0,
		}
	}
	var lastId LastId
	err = yaml.Unmarshal(data, &lastId)
	if err != nil {
		fmt.Println("yaml.Unmarshal: ", err)
		return LastId{
			World: 0,
			Area: 0,
			Person: 0,
			Gender: 0,
			Race: 0,
			Job: 0,
		}
	}
	return lastId
}

