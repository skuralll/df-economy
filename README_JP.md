# DF Economy

Minecraft Bedrock Edition用の経済システムライブラリ（**df-mc/dragonfly**フレームワーク対応）

## 概要

DF Economyは、Minecraftサーバー内で使用できる柔軟な通貨システムを提供します。マルチデータベース対応で、プレイヤー間での送金、残高確認、ランキング表示などの機能を備えています。SQLite、MySQL、PostgreSQLに対応しています。

## ゲーム内コマンド

| コマンド | 説明 | 使用例 |
| --- | --- | --- |
| `/economy` | コマンドヘルプを表示 | `/economy` |
| `/economy balance [プレイヤー名]` | 残高を表示 | `/economy balance` または `/economy balance Steve` |
| `/economy pay <プレイヤー名> <金額>` | 他のプレイヤーに送金 | `/economy pay Steve 100` |
| `/economy set <プレイヤー名> <金額>` | 残高を設定（設定可能） | `/economy set Steve 1000` |
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
    "context"
    "github.com/df-mc/dragonfly/server"
    "github.com/skuralll/df-permission/permission"
    "github.com/skuralll/dfeconomy/dragonfly/commands"
    "github.com/skuralll/dfeconomy/economy/config"
    "github.com/skuralll/dfeconomy/economy/service"
)

func main() {
    // 設定付きで経済システムサービスを作成
    cfg := config.Config{
        DBType:         "sqlite",          // または "mysql", "postgres"
        DBDSN:          "./economy.db",    // データベース接続文字列
        DefaultBalance: 100.0,             // 新規プレイヤーの初期残高
        EnableSetCmd:   false,             // /economy setコマンドを有効化 (デフォルト: false)
    }
    
    // 権限管理システムを初期化
    // 権限管理システムが不要な場合は nilを渡すことも可能
    pMgr := permission.NewManager()
    
    svc, cleanup, err := service.NewEconomyService(cfg, pMgr)
    if err != nil {
        panic(err)
    }
    defer cleanup()
    
    // コマンドを登録
    commands.RegisterCommands(svc, cfg)
    
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

経済システムは複数のデータベースバックエンドをサポートしています：

#### SQLite（デフォルト）
```go
cfg := config.Config{
    DBType: "sqlite",
    DBDSN:  "./economy.db",
    DefaultBalance: 100.0,
    EnableSetCmd: false,
}
```

#### MySQL
```go
cfg := config.Config{
    DBType: "mysql",
    DBDSN:  "user:password@tcp(localhost:3306)/economy?charset=utf8mb4&parseTime=True&loc=Local",
    DefaultBalance: 100.0,
    EnableSetCmd: false,
}
```

#### PostgreSQL
```go
cfg := config.Config{
    DBType: "postgres",
    DBDSN:  "host=localhost user=user password=password dbname=economy port=5432 sslmode=disable",
    DefaultBalance: 100.0,
    EnableSetCmd: false,
}
```

データベースのテーブルとスキーマは起動時に自動作成されます。

#### setコマンドの有効化（オプション）
残高管理用の`/economy set`コマンドを有効化する場合：
```go
cfg := config.Config{
    DBType: "sqlite",
    DBDSN:  "./economy.db",
    DefaultBalance: 100.0,
    EnableSetCmd: true,  // setコマンドを有効化
}
```

**注意**: setコマンドはセキュリティ上の理由からデフォルトで無効化されています。

## 機能

- **マルチデータベース対応**: SQLite、MySQL、PostgreSQLをサポート
- **残高管理**: プレイヤーの残高確認と設定
- **送金システム**: プレイヤー間での安全な送金
- **ランキング**: 残高によるプレイヤーランキング
- **自動登録**: 設定可能な初期残高での新規プレイヤー自動登録
- **エラーハンドリング**: 適切な検証付きでわかりやすいエラーメッセージ
- **CGO不要**: 全データベースドライバーのPure Go実装
- **トランザクション安全性**: 適切なロールバック処理付きのACID準拠
- **コマンド制御**: セキュリティ強化のための設定可能なコマンド有効性

## 要件

- Go 1.24以上
- df-mc/dragonfly フレームワーク
- データベース: SQLite（自動セットアップ）、MySQL、またはPostgreSQL

## ライセンス

MIT License