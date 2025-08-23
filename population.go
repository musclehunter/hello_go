package main

import (
    "github.com/brianvoe/gofakeit/v6"
)

// initResidents は各エリアに initial_population 人の住人を追加します。
// ただし各エリアの人口は max_population を上限とします。
func initResidents(cfg Config) {
    worldMu.Lock()
    defer worldMu.Unlock()

    for ai := range World.Areas {
        area := &World.Areas[ai]
        // 既存の人口を考慮して追加人数を決定
        remaining := cfg.MaxPopulation - area.Population
        if remaining <= 0 {
            continue
        }
        toAdd := cfg.InitialPopulation
        if toAdd > remaining {
            toAdd = remaining
        }
        for i := 0; i < toAdd; i++ {
            p := NewPerson(gofakeit.Name())
            // World に登録
            World.Persons = append(World.Persons, p)
            // Area にIDを登録
            area.ResidentIDs = append(area.ResidentIDs, p.Id)
            area.Population++
            // 追加時にその人の収入を反映 (Level * Job.Income)
            if job, ok := World.GetJobById(p.MainJobID); ok {
                area.Income += p.Level * job.Income
            }
        }
    }
}

// addResidentsPerTurn は各ターンで各エリアに最大5人の住人を追加します。
// ただし各エリアの人口は MaxPopulation を上限とします。
// 追加した住人は直ちに行動ゴルーチンを開始します。
func addResidentsPerTurn() {
    worldMu.Lock()
    defer worldMu.Unlock()

    for ai := range World.Areas {
        area := &World.Areas[ai]
        remaining := appConfig.MaxPopulation - area.Population
        if remaining <= 0 {
            continue
        }
        toAdd := 5
        if toAdd > remaining {
            toAdd = remaining
        }
        for i := 0; i < toAdd; i++ {
            p := NewPerson(gofakeit.Name())
            // World に登録
            World.Persons = append(World.Persons, p)
            // Area にIDを登録
            area.ResidentIDs = append(area.ResidentIDs, p.Id)
            area.Population++
            if job, ok := World.GetJobById(p.MainJobID); ok {
                area.Income += p.Level * job.Income
            }
            // 住人の行動ゴルーチンを開始
            wg.Add(1)
            go p.Run(appCtx, area.Id, events, &wg)
        }
    }
}
