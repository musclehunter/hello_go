package main

type Job struct {
    Id   int
    Name string
    Income int
    // DeathRate は 0.0 - 1.0 の確率で行動時に死亡する確率
    DeathRate float64
}

func NewJob(name string, income int, deathRate float64) Job {
    LastId.Job++
    return Job{
        Id:   LastId.Job,
        Name: name,
        Income: income,
        DeathRate: deathRate,
    }
}
