package main

type Race struct {
	Id   int
	Name string
}

func NewRace(name string) Race {
	LastId.Race++
	return Race{
		Id:   LastId.Race,
		Name: name,
	}
}
