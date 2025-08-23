package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"fmt"
	"path/filepath"
)

type LastIdStruct struct {
	World int
	Area int
	Person int
	Gender int
	Race int
	Job int
}

func SaveLastId() {
	if err := ensureDataDir(); err != nil {
		fmt.Println("ensureDataDir:", err)
		return
	}
	data, err := yaml.Marshal(LastId)
	if err != nil {
		fmt.Println("yaml.Marshal: ", err)
		return
	}
	ioutil.WriteFile(filepath.Join(dataDir, "last_id.yml"), data, 0644)
}

func loadLastId() LastIdStruct {
	if err := ensureDataDir(); err != nil {
		fmt.Println("ensureDataDir:", err)
		return LastIdStruct{
			World: 0,
			Area: 0,
			Person: 0,
			Gender: 0,
			Race: 0,
			Job: 0,
		}
	}
	data, err := ioutil.ReadFile(filepath.Join(dataDir, "last_id.yml"))
	if err != nil {
		fmt.Println("ioutil.ReadFile: ", err)
		return LastIdStruct{
			World: 0,
			Area: 0,
			Person: 0,
			Gender: 0,
			Race: 0,
			Job: 0,
		}
	}
	var lastId LastIdStruct
	err = yaml.Unmarshal(data, &lastId)
	if err != nil {
		fmt.Println("yaml.Unmarshal: ", err)
		return LastIdStruct{
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

var LastId LastIdStruct = loadLastId()