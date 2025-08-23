package main

import (
    "context"
    "sync"
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
)
