package main

import (
	"fmt"
	"os"
	"time"
	"net/http"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	"context"
	"strconv"
)

var Turn int = 0

func main() {
	gofakeit.Seed(time.Now().UnixNano())
	config, err := LoadConfig("config.yml")
	if err != nil {
		fmt.Println("config.yml not found", err)
		os.Exit(1)
	}
	// 設定をグローバルに保持
	appConfig = config

	// ランタイム初期化と初期住人生成
	appCtx, cancel = context.WithCancel(context.Background())
	events = make(chan Event, 4096)
	initResidents(config)

	// 住人ゴルーチン起動（読み取りロックで安全に参照）
	worldMu.RLock()
	for _, area := range World.Areas {
		for _, pid := range area.ResidentIDs {
			if idx, ok := World.GetPersonIndexById(pid); ok {
				p := World.Persons[idx]
				wg.Add(1)
				go p.Run(appCtx, area.Id, events, &wg)
			}
		}
	}
	worldMu.RUnlock()

	// 単一適用ループ
	go applyEvents()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		worldMu.RLock()
		defer worldMu.RUnlock()
		json.NewEncoder(w).Encode(World)
	})

	// アプリ内ログ取得エンドポイント
	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		n := 100
		if q := r.URL.Query().Get("n"); q != "" {
			if v, err := strconv.Atoi(q); err == nil && v > 0 {
				n = v
			}
		}
		json.NewEncoder(w).Encode(getLogs(n))
	})

	go func() {
		addr := fmt.Sprintf(":%d", config.Port)
		fmt.Printf("listening on %s\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Println("ListenAndServe: ", err)
		}
	}()

	last := time.Now()
	for {
		now := time.Now()
		if now.Sub(last) >= time.Duration(config.TurnSeconds) * time.Second {
			action()
			last = now
			Turn++
			if (Turn >= config.MaxTurn) {
				Shutdown()
			}
		}
	}
}

func action() {
	// 各ターンで各エリアに最大5人追加（MaxPopulationを超えない）
	addResidentsPerTurn()
	// ファイル保存のみ（標準出力は出さない）
	SaveWorld()
	SaveLastId()
	time.Sleep(100 * time.Millisecond)
}