package main

import (
    "math/rand"
	"fmt"
    "time"
)

// Event は住人から送られるイベントの共通インタフェース
// 単一の適用ゴルーチンが受け取り、World へ状態変更を適用します。
type Event interface{}

// PersonHeartbeat は住人が周期的に行動したことを示す簡易イベント
// 実際のアクション種別が増えたら型を追加してください。
type PersonHeartbeat struct {
    PersonID int
    AreaID   int
}

// PersonActed は住人の1ターンの行動終了を示すイベント
// 単一適用ループで Level++ と Income 増分を適用します。
type PersonActed struct {
    PersonID int
    AreaID   int
}

// applyEvents は単一ゴルーチンでイベントを受け取り、World へ反映する
func applyEvents() {
    for {
        select {
        case ev, ok := <-events:
            if !ok {
                return
            }
            switch v := ev.(type) {
            case PersonHeartbeat:
                // 今はダミー。必要なら最終アクティブ時刻などを管理
                _ = v
            case PersonActed:
                // レベルアップと収入加算を適用（正規化構造）
                worldMu.Lock()
                // Area を特定
                for ai := range World.Areas {
                    if World.Areas[ai].Id != v.AreaID {
                        continue
                    }
                    // Person を特定（World.Personsから）
                    if pi, ok := World.GetPersonIndexById(v.PersonID); ok {
                        p := World.Persons[pi]
                        // すでに死亡しているなら無視
                        if !p.Alive {
                            break
                        }
                        // 死亡判定（MainJobIDのJobからDeathRate参照）
                        job, _ := World.GetJobById(p.MainJobID)
                        if rand.Float64() < job.DeathRate {
                            // 収入からこの住人分を減算
                            World.Areas[ai].Income -= p.Level * job.Income
                            if World.Areas[ai].Income < 0 {
                                World.Areas[ai].Income = 0
                            }
                            // AreaのResidentIDsから削除
                            ids := World.Areas[ai].ResidentIDs
                            for idx, id := range ids {
                                if id == p.Id {
                                    World.Areas[ai].ResidentIDs = append(ids[:idx], ids[idx+1:]...)
                                    break
                                }
                            }
                            // 人口減少
                            if World.Areas[ai].Population > 0 {
                                World.Areas[ai].Population--
                            }
                            // フラグ更新（互換用、すぐ削除するが念のため）
                            p.Alive = false
                            p.DeathDay = time.Now()
                            // DeadPersons に退避してから Persons から削除
                            World.DeadPersons = append(World.DeadPersons, p)
                            World.Persons = append(World.Persons[:pi], World.Persons[pi+1:]...)
                            // 死亡ログ
                            fmt.Println(p.Name, "(Level ", p.Level, ") ", " died in Area ", World.Areas[ai].Name)
                            // Persons自体は履歴保持のため残す（完全削除したい場合はここで削除）
                        } else {
                            // 生存: レベルアップと収入増分
                            p.Level++
                            World.Persons[pi] = p
                            World.Areas[ai].Income += job.Income
                        }
                    }
                    break
                }
                worldMu.Unlock()
            default:
                // 未知イベント: 何もしない
            }
        }
    }
}
