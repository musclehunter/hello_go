package main

import (
    "context"
    "sync"
    "time"
    "strconv"
)

var (
    appCtx context.Context
    cancel context.CancelFunc
    wg     sync.WaitGroup

    events chan Event

    // World への読み書きの同期用
    worldMu sync.RWMutex

    // アプリ設定（MaxPopulation 等を参照）
    appConfig Config

    // アプリ内ログ（メモリリングバッファ）
    logs   []AppLog
    logsMu sync.RWMutex
)

// 1件のアプリ内ログ
type AppLog struct {
    Time    time.Time `json:"time"`
    Type    string    `json:"type"`
    Message string    `json:"message"`
}

const maxLogs = 1000

// ログを追加
func addLog(typ, msg string) {
    logsMu.Lock()
    defer logsMu.Unlock()
    if len(logs) >= maxLogs {
        // 先頭を落としてリング的に保持
        logs = logs[1:]
    }
    logs = append(logs, AppLog{Time: time.Now(), Type: typ, Message: msg})
}

// 末尾から最大n件を取得（新しい順）
func getLogs(n int) []AppLog {
    logsMu.RLock()
    defer logsMu.RUnlock()
    if n <= 0 || n > len(logs) {
        n = len(logs)
    }
    out := make([]AppLog, n)
    copy(out, logs[len(logs)-n:])
    return out
}

// int を文字列化（ログ用の簡易ヘルパ）
func fmtInt(i int) string { return strconv.Itoa(i) }
