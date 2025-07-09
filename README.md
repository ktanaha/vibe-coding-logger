# 🎯 Vibe Coding Logger

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/ktanaha/vibe-coding-logger/workflows/CI/badge.svg)](https://github.com/ktanaha/vibe-coding-logger/actions)

バイブコーディング（思考プロセスを重視したプログラミング）に特化した強力なGoロギングライブラリです。

## ✨ 特徴

- 🧠 **バイブコーディング専用機能**: 思考プロセス、決定、学習、ブレイクスルーなどを詳細に記録
- 📊 **操作追跡**: 入力・出力・処理内容を自動的に追跡
- 🖥️ **自動システム情報収集**: OS、言語バージョン、Git情報、環境変数を自動記録
- 🔄 **高度なエラーハンドリング**: リトライ、リカバリー、サーキットブレーカー機能
- 📋 **構造化ログ**: JSON、テキスト形式での柔軟な出力
- 🔍 **トレーシング対応**: 分散システムでの追跡に対応
- ⚡ **高パフォーマンス**: 並行処理とバッファリング、キャッシュによる最適化

## 🚀 クイックスタート

### インストール

```bash
go get github.com/ktanaha/vibe-coding-logger
```

### 基本的な使用方法

```go
package main

import "github.com/ktanaha/vibe-coding-logger/pkg/logger"

func main() {
    // デフォルトロガーを作成
    log := logger.Default()
    
    // 基本的なログ出力
    log.Info("アプリケーション開始")
    log.Error("エラーが発生", logger.String("component", "database"))
    
    // フィールド付きログ
    log.Info("ユーザーログイン", 
        logger.String("user_id", "user123"),
        logger.Int("attempt", 1))
}
```

### バイブコーディング機能

```go
// バイブトラッカーの作成
vibeTracker := logger.NewVibeTracker(log, "session_123", "web_api", "coding")

// 思考プロセスの記録
vibeTracker.LogThinkingProcess(
    "認証方式の検討",
    []string{"JWT", "Session", "OAuth2"})

// 決定の記録
vibeTracker.LogDecision(
    "JWT認証を採用",
    "ステートレスで拡張性が高い",
    []string{"Session認証", "OAuth2"})

// コード変更の記録
vibeTracker.LogCodeChange(
    "auth.go",
    "新規作成",
    "",
    "func AuthHandler() { ... }",
    "JWT認証ハンドラーの実装")
```

## 📖 詳細ドキュメント

### 操作追跡

```go
// 操作の開始
tracker := log.StartOperation("user_registration", map[string]interface{}{
    "email": "user@example.com",
    "username": "newuser",
})

// 処理...
time.Sleep(100 * time.Millisecond)

// 操作の完了
log.CompleteOperation(tracker, map[string]interface{}{
    "user_id": "user123",
    "status": "active",
})
```

### システム情報の自動記録

```go
// システム情報の有効化/無効化
log.EnableSystemInfo(true)   // デフォルト: 有効
log.EnableRuntimeInfo(false) // デフォルト: 無効（パフォーマンス考慮）

// 現在の設定確認
fmt.Println("システム情報:", log.IsSystemInfoEnabled())
fmt.Println("ランタイム情報:", log.IsRuntimeInfoEnabled())
```

#### 収集される情報

**システム情報:**
- OS（Linux, Windows, macOS）
- アーキテクチャ（amd64, arm64等）
- Go言語バージョン
- CPU数、ホスト名、プロセスID

**環境情報:**
- 作業ディレクトリ、GOPATH、GOROOT
- Git ブランチ、コミットハッシュ、リポジトリURL
- エディタ情報
- Node.js、Python、Dockerバージョン

**ランタイム情報:**
- Goroutine数、メモリ使用量
- GC統計、スタック使用量

### エラーハンドリング

```go
// バイブエラーハンドラー
errorHandler := logger.NewVibeErrorHandler(log, "session_123", "web_api")

// コーディングエラーの処理
errorHandler.HandleCodingError(
    err,
    "models/user.go",
    45,
    "db.Query(sql, params...)",
    "接続プールの設定を見直し")

// リトライハンドラー
retryHandler := logger.NewRetryHandler(log)
err := retryHandler.ExecuteWithRetry("api_call", func() error {
    return callExternalAPI()
}, 3, 1*time.Second, map[string]interface{}{
    "endpoint": "https://api.example.com",
})
```

### カスタム設定

```go
// カスタムロガーの作成
log := logger.New(logger.DEBUG)

// ファイルライターの追加
fileWriter, _ := logger.NewFileWriter("app.log")
log.AddWriter(fileWriter)

// JSONフォーマッターの設定
log.SetFormatter(logger.NewJSONFormatter())

// ローテーション設定
rotatingWriter, _ := logger.NewRotatingFileWriter("app.log", 10*1024*1024, 5)
log.AddWriter(rotatingWriter)
```

## 🎯 バイブコーディングでの活用

### 推奨ワークフロー

1. **セッション開始時**: 問題の理解と目標設定、環境情報を記録
2. **思考プロセス**: 検討した選択肢と思考の流れを記録
3. **決定の記録**: なぜその選択をしたかの理由を明確に
4. **コード変更**: 変更の理由と期待する効果を記録
5. **テスト結果**: 成功/失敗とその原因を詳細に
6. **学習内容**: 新しく理解したことや気づきを記録
7. **セッション終了**: 成果と次のステップ、環境変化を整理

### システム情報の活用場面

- **環境の再現**: 問題が発生した時の正確な環境情報
- **パフォーマンス分析**: ハードウェア構成とパフォーマンスの関係
- **バージョン管理**: 言語やツールのバージョンによる動作の違い
- **チーム共有**: 他の開発者と環境情報を共有

## ⚡ パフォーマンス最適化

```go
// バッファリングライターの使用
bufferedWriter, _ := logger.NewBufferedFileWriter("app.log", 100)
log.AddWriter(bufferedWriter)

// ログレベルの適切な設定
log.SetLevel(logger.INFO) // 本番環境
log.SetLevel(logger.DEBUG) // 開発環境

// システム情報の選択的有効化
log.EnableSystemInfo(true)   // 基本システム情報（軽量）
log.EnableRuntimeInfo(false) // ランタイム情報（重い、デバッグ時のみ）

// 本番環境では最小限に
if isProduction {
    log.EnableSystemInfo(false)
    log.EnableRuntimeInfo(false)
    log.SetLevel(logger.WARN)
}
```

## 📁 プロジェクト構造

```
vibe-coding-logger/
├── pkg/logger/              # 公開API
│   ├── interfaces.go        # インターフェース定義
│   ├── logger.go           # メインロガー実装
│   ├── system_info.go      # システム情報収集
│   ├── tracker.go          # 操作・バイブトラッカー
│   └── error_handler.go    # エラーハンドリング
├── internal/               # 内部実装
│   ├── formatter/          # ログフォーマッター
│   └── writer/             # ログライター
├── examples/               # 使用例
├── tests/                  # テストコード
├── docs/                   # ドキュメント
└── README.md
```

## 🤝 コントリビューション

コントリビューションを歓迎します！

1. このリポジトリをフォーク
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

詳細は [CONTRIBUTING.md](CONTRIBUTING.md) をご覧ください。

## 📄 ライセンス

このプロジェクトは [MIT License](LICENSE) の下で公開されています。

## 🙏 謝辞

- バイブコーディング概念の発案者
- Goコミュニティからの貴重なフィードバック
- オープンソースライブラリのメンテナー

## 📧 サポート

- **Issues**: [GitHub Issues](https://github.com/ktanaha/vibe-coding-logger/issues)
- **Discussions**: [GitHub Discussions](https://github.com/ktanaha/vibe-coding-logger/discussions)
- **Documentation**: [Wiki](https://github.com/ktanaha/vibe-coding-logger/wiki)

---

<div align="center">

**バイブコーディングの体験向上に向けて、一緒に開発しましょう！**

[⭐ Star this repository](https://github.com/ktanaha/vibe-coding-logger) | [📖 Read the docs](https://github.com/ktanaha/vibe-coding-logger/wiki) | [🐛 Report a bug](https://github.com/ktanaha/vibe-coding-logger/issues)

</div>