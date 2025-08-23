package main

type Gender struct {
	Id   int
	Name string
}

func NewGender(name string) Gender {
	LastId.Gender++
	return Gender{
		Id:   LastId.Gender,
		Name: name,
	}
}
