package main

import (
	"time"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"os"
	"path/filepath"
)

type WorldStruct struct {
	Id int
	Name string
	BirthDay time.Time
	Areas []Area
	Genders []Gender
	Races []Race
	Jobs []Job
	Persons []Person
	DeadPersons []Person
}

// IDからGenderを取得
func (w *WorldStruct) GetGenderById(id int) (Gender, bool) {
	for _, g := range w.Genders {
		if g.Id == id {
			return g, true
		}
	}
	return Gender{}, false
}

// IDからRaceを取得
func (w *WorldStruct) GetRaceById(id int) (Race, bool) {
	for _, r := range w.Races {
		if r.Id == id {
			return r, true
		}
	}
	return Race{}, false
}

// IDからJobを取得
func (w *WorldStruct) GetJobById(id int) (Job, bool) {
	for _, j := range w.Jobs {
		if j.Id == id {
			return j, true
		}
	}
	return Job{}, false
}

// IDからPersonのインデックスを取得
func (w *WorldStruct) GetPersonIndexById(id int) (int, bool) {
	for i, p := range w.Persons {
		if p.Id == id {
			return i, true
		}
	}
	return -1, false
}

func NewWorld(name string) WorldStruct {
	LastId.World++
	return WorldStruct{
		Id: LastId.World,
		Name: name,
		BirthDay: time.Now(),
		Areas: []Area{},
		Genders: []Gender{},
		Races: []Race{},
		Jobs: []Job{},
		Persons: []Person{},
		DeadPersons: []Person{},
	}
}

const dataDir = "data"

func ensureDataDir() error {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return os.MkdirAll(dataDir, 0755)
	}
	return nil
}

// Worldを分割YAMLに保存する（正規化済み構造をそのまま保存）
func SaveWorld() {
	if err := ensureDataDir(); err != nil {
		fmt.Println("ensureDataDir:", err)
		return
	}

	// world meta
	wm := struct{ Id int; Name string; BirthDay time.Time }{Id: World.Id, Name: World.Name, BirthDay: World.BirthDay}
	if data, err := yaml.Marshal(wm); err == nil {
		_ = ioutil.WriteFile(filepath.Join(dataDir, "world.yml"), data, 0644)
	} else {
		fmt.Println("yaml.Marshal world:", err)
	}

	// genders
	if data, err := yaml.Marshal(World.Genders); err == nil {
		_ = ioutil.WriteFile(filepath.Join(dataDir, "genders.yml"), data, 0644)
	} else {
		fmt.Println("yaml.Marshal genders:", err)
	}

	// races
	if data, err := yaml.Marshal(World.Races); err == nil {
		_ = ioutil.WriteFile(filepath.Join(dataDir, "races.yml"), data, 0644)
	} else {
		fmt.Println("yaml.Marshal races:", err)
	}

	// jobs
	if data, err := yaml.Marshal(World.Jobs); err == nil {
		_ = ioutil.WriteFile(filepath.Join(dataDir, "jobs.yml"), data, 0644)
	} else {
		fmt.Println("yaml.Marshal jobs:", err)
	}

	// areas（ResidentIDsを含む正規化構造のまま保存）
	if data, err := yaml.Marshal(World.Areas); err == nil {
		_ = ioutil.WriteFile(filepath.Join(dataDir, "areas.yml"), data, 0644)
	} else {
		fmt.Println("yaml.Marshal areas:", err)
	}

    // persons（生存者のみ保存）
    var alive []Person
    for _, p := range World.Persons {
        if p.Alive {
            alive = append(alive, p)
        }
    }
    if data, err := yaml.Marshal(alive); err == nil {
        _ = ioutil.WriteFile(filepath.Join(dataDir, "persons_alive.yml"), data, 0644)
    } else {
        fmt.Println("yaml.Marshal persons_alive:", err)
    }
    // persons_dead（死亡者を保存）
    if data, err := yaml.Marshal(World.DeadPersons); err == nil {
        _ = ioutil.WriteFile(filepath.Join(dataDir, "persons_dead.yml"), data, 0644)
    } else {
        fmt.Println("yaml.Marshal persons_dead:", err)
    }
}

func loadWorld() WorldStruct {
	// world meta
	metaPath := filepath.Join(dataDir, "world.yml")
	data, err := ioutil.ReadFile(metaPath)
	if err != nil {
		// 初回起動など: 新規作成（レガシーは無視して破棄）
		fmt.Println("read world meta:", err)
		return createWorld()
	}
	var wm struct{ Id int; Name string; BirthDay time.Time }
	if err := yaml.Unmarshal(data, &wm); err != nil {
		fmt.Println("unmarshal world meta:", err)
		return createWorld()
	}

	// genders
	gendersPath := filepath.Join(dataDir, "genders.yml")
	data, err = ioutil.ReadFile(gendersPath)
	if err != nil {
		fmt.Println("read genders:", err)
		return createWorld()
	}
	var genders []Gender
	if err := yaml.Unmarshal(data, &genders); err != nil {
		fmt.Println("unmarshal genders:", err)
		return createWorld()
	}

	// races
	racesPath := filepath.Join(dataDir, "races.yml")
	data, err = ioutil.ReadFile(racesPath)
	if err != nil {
		fmt.Println("read races:", err)
		return createWorld()
	}
	var races []Race
	if err := yaml.Unmarshal(data, &races); err != nil {
		fmt.Println("unmarshal races:", err)
		return createWorld()
	}

	// jobs
	jobsPath := filepath.Join(dataDir, "jobs.yml")
	data, err = ioutil.ReadFile(jobsPath)
	if err != nil {
		fmt.Println("read jobs:", err)
		return createWorld()
	}
	var jobs []Job
	if err := yaml.Unmarshal(data, &jobs); err != nil {
		fmt.Println("unmarshal jobs:", err)
		return createWorld()
	}

	// areas（ResidentIDsを含む）
	areasPath := filepath.Join(dataDir, "areas.yml")
	data, err = ioutil.ReadFile(areasPath)
	if err != nil {
		fmt.Println("read areas:", err)
		return createWorld()
	}
	var areas []Area
	if err := yaml.Unmarshal(data, &areas); err != nil {
		fmt.Println("unmarshal areas:", err)
		return createWorld()
	}

    // persons（ID参照型のまま読込）: persons_alive.yml を優先、無ければ persons.yml にフォールバック
    var persons []Person
    alivePath := filepath.Join(dataDir, "persons_alive.yml")
    aliveData, errAlive := ioutil.ReadFile(alivePath)
    if errAlive == nil {
        var alivePersons []Person
        if err := yaml.Unmarshal(aliveData, &alivePersons); err == nil {
            persons = append(persons, alivePersons...)
        } else {
            fmt.Println("unmarshal persons_alive:", err)
        }
    } else {
        // persons_alive.yml が無ければ空のまま
    }

    // レガシーデータ互換: Aliveフィールドが無い旧データを読み込んだ場合（ゼロ値）
    for i := range persons {
        if !persons[i].Alive && persons[i].DeathDay.IsZero() {
            persons[i].Alive = true
        }
    }
    world := WorldStruct{Id: wm.Id, Name: wm.Name, BirthDay: wm.BirthDay, Areas: areas, Genders: genders, Races: races, Jobs: jobs, Persons: persons}
    // 読み込み後に派生値を再計算（存在しないIDは除去）
    for i := range world.Areas {
        area := &world.Areas[i]
        // ResidentIDs を生存者に限定し、存在しないIDは除去
        filtered := make([]int, 0, len(area.ResidentIDs))
        income := 0
        for _, pid := range area.ResidentIDs {
            if pi, ok := world.GetPersonIndexById(pid); ok {
                p := world.Persons[pi]
                if p.Alive { // 念のため
                    filtered = append(filtered, pid)
                    if job, ok := world.GetJobById(p.MainJobID); ok {
                        income += p.Level * job.Income
                    }
                }
            }
        }
        area.ResidentIDs = filtered
        // 人口はResidentIDsの長さ
        area.Population = len(area.ResidentIDs)
        // 収入は resident の (Level * Job.Income) の合計
        area.Income = income
    }
    return world
}

func createWorld() WorldStruct {
	world := NewWorld(gofakeit.Fruit() + " World")
	world.Areas = append(world.Areas, NewArea(gofakeit.BeerName() + " Area"))
	world.Areas = append(world.Areas, NewArea(gofakeit.BeerName() + " Area"))
	world.Areas = append(world.Areas, NewArea(gofakeit.BeerName() + " Area"))
	world.Genders = append(world.Genders, NewGender("Male"))
	world.Genders = append(world.Genders, NewGender("Female"))
	world.Races = append(world.Races, NewRace("Human"))
	world.Races = append(world.Races, NewRace("Elf"))
	world.Races = append(world.Races, NewRace("Orc"))
	world.Jobs = append(world.Jobs, NewJob("Farmer", 10, 0.01))
	world.Jobs = append(world.Jobs, NewJob("Merchant", 20, 0.05))
	world.Jobs = append(world.Jobs, NewJob("Soldier", 30, 0.1))
	return world
}

var World WorldStruct = loadWorld()