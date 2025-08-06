package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"net/http"
	"encoding/json"
)

func main() {
	numbers := []int{1, 1, 1, 1, 1}

	config, err := LoadConfig("config.yml")
	if err != nil {
		fmt.Println("config.yml not found", err)
		os.Exit(1)
	}

	http.HandleFunc("/numbers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(numbers)
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
			action(numbers)
			last = now
			count++
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func action(numbers []int) {
	// ランダムにどれかをインクリメントする
	numbers[rand.Intn(len(numbers))]++

	fmt.Printf("[%s] %v\n", time.Now().Format("2006-01-02 15:04:05"), numbers)
	os.Stdout.Sync()
}