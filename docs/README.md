# DF Economy

Minecraft Bedrock Edition用の経済システムライブラリ（**df-mc/dragonfly**フレームワーク対応）

## 概要

DF Economyは、Minecraftサーバー内で使用できる通貨システムを提供します。プレイヤー間での送金、残高確認、ランキング表示などの機能を備えています。

## ゲーム内コマンド

| コマンド | 説明 | 使用例 |
| --- | --- | --- |
| `/economy` | コマンドヘルプを表示 | `/economy` |
| `/economy balance [プレイヤー名]` | 残高を表示 | `/economy balance` または `/economy balance Steve` |
| `/economy pay <プレイヤー名> <金額>` | 他のプレイヤーに送金 | `/economy pay Steve 100` |
| `/economy set <プレイヤー名> <金額>` | 残高を設定（管理者用） | `/economy set Steve 1000` |
| `/economy top <ページ>` | 残高ランキングを表示 | `/economy top 1` |

## 使用方法

### 1. ライブラリのインストール

```bash
go get github.com/skuralll/dfeconomy
```

### 2. サーバーへの組み込み

```go
package main

import (
    "github.com/df-mc/dragonfly/server"
    "github.com/skuralll/dfeconomy/dragonfly/commands"
    "github.com/skuralll/dfeconomy/economy/service"
)

func main() {
    // 経済システムサービスを作成
    svc, cleanup, err := service.NewEconomyService()
    if err != nil {
        panic(err)
    }
    defer cleanup()
    
    // コマンドを登録
    commands.RegisterCommands(svc)
    
    // サーバーの設定とスタート
    srv := server.DefaultConfig().New()
    srv.Listen()
    
    // プレイヤーの参加時に自動登録
    for p := range srv.Accept() {
        svc.RegisterUser(context.Background(), p.UUID(), p.Name())
    }
}
```

### 3. 設定

データベースファイル（`foo.db`）が自動的に作成されます。プレイヤーの残高情報はここに保存されます。

## 機能

- **残高管理**: プレイヤーの残高確認と設定
- **送金システム**: プレイヤー間での安全な送金
- **ランキング**: 残高によるプレイヤーランキング
- **自動登録**: 新規プレイヤーの自動登録
- **エラーハンドリング**: わかりやすいエラーメッセージ

## 要件

- Go 1.21以上
- df-mc/dragonfly フレームワーク
- SQLite（自動セットアップ）

## ライセンス

MIT License