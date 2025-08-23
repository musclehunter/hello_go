package main

import (
    "context"
    "time"
)

// Shutdown は全住人ゴルーチンの停止を指示し、完了を待ち、イベントチャネルを閉じる
func Shutdown() {
    if cancel != nil {
        cancel()
    }
    // 念のため待機にタイムアウトを設けたい場合は Context を使う
    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()

    select {
    case <-done:
    case <-time.After(5 * time.Second):
        // タイムアウト（ログ等を入れてもよい）
    }

    if events != nil {
        // applyEvents の select で close を検知して終了
        close(events)
    }
    // サーバ優雅停止等を行いたい場合は http.Server を使う実装に変更
    _ = context.Canceled
}
