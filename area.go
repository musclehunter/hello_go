package main


type Area struct {
    Id         int
    Name       string
    Population int
    ResidentIDs  []int
    Income    int
}

func NewArea(name string) Area {
    LastId.Area++
    return Area{
        Id:         LastId.Area,
        Name:       name,
        Population: 0,
        ResidentIDs:  []int{},
        Income:    0,
    }
}
