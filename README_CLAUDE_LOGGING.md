# Claude Code ログ出力機能

このプロジェクトには、Claude Code実行中のログをトップディレクトリにテキストファイルとして出力する機能が含まれています。

## 機能概要

- Claude Codeセッション中のすべての活動をログファイルに記録
- 日時付きファイル名での自動ログファイル生成
- 構造化されたログ出力（タスク、コマンド、ファイル操作、エラーなど）
- 日本語対応のログメッセージ

## 使用方法

### 1. 基本的な使用方法

```go
package main

import (
    "log"
)

func main() {
    // Claude Codeロガーを初期化
    ccLogger, err := NewClaudeCodeLogger()
    if err != nil {
        log.Fatalf("ロガーの初期化に失敗しました: %v", err)
    }
    defer ccLogger.Close()

    // ログファイルパスを表示
    fmt.Printf("ログファイル: %s\n", ccLogger.GetLogFilePath())
    
    // 各種ログを記録
    ccLogger.LogTaskStart("ファイル作成", "新しいファイルを作成します")
    ccLogger.LogFileOperation("CREATE", "example.go")
    ccLogger.LogTaskComplete("ファイル作成", time.Second*2)
}
```

### 2. サンプル実行

```bash
# ログ出力のサンプルを実行
go run claude_code_logger.go
```

実行すると、以下のようなファイルが作成されます：
- `claude_code_log_20240119_143025.txt`

## ログファイル形式

ログファイルはテキスト形式で以下の情報を含みます：

```
[2024-01-19 14:30:25] INFO Claude Code セッション開始 log_file=./claude_code_log_20240119_143025.txt start_time="2024-01-19 14:30:25"

[2024-01-19 14:30:25] INFO タスク開始 task_name="ファイル作成" description="新しいGoファイルを作成します" timestamp="2024-01-19 14:30:25"

[2024-01-19 14:30:25] INFO ファイル操作 operation=CREATE file_path=example.go timestamp="2024-01-19 14:30:25"

[2024-01-19 14:30:25] INFO コマンド実行 command=go arguments=[build .] timestamp="2024-01-19 14:30:25"

[2024-01-19 14:30:25] INFO ユーザー対話 user_message="Goファイルを作成してください" claude_response="新しいGoファイルを作成しました。ビルドも正常に完了しました。" timestamp="2024-01-19 14:30:25"

[2024-01-19 14:30:25] INFO ツール使用 tool_name=Write purpose="ファイル作成" parameters=map[content_length:150 file_path:example.go] timestamp="2024-01-19 14:30:25"

[2024-01-19 14:30:25] INFO コード生成 file_name=example.go language=Go description="サンプルのGo言語プログラム" timestamp="2024-01-19 14:30:25"

[2024-01-19 14:30:25] INFO タスク完了 task_name="ファイル作成" duration=2s timestamp="2024-01-19 14:30:25"

[2024-01-19 14:30:25] DEBUG デバッグ情報 context="プロセス状況" details=map[active_goroutines:5 cpu_usage:2.5% memory_usage:15MB] timestamp="2024-01-19 14:30:25"

[2024-01-19 14:30:25] INFO Claude Code セッション終了 end_time="2024-01-19 14:30:25"
```

## 利用可能なログ機能

### 基本ログ機能
- `LogClaudeCodeStart()` - Claude Codeセッション開始
- `LogCommand(command, args)` - コマンド実行
- `LogFileOperation(operation, filePath)` - ファイル操作
- `LogError(operation, err)` - エラー記録

### タスク管理
- `LogTaskStart(taskName, description)` - タスク開始
- `LogTaskComplete(taskName, duration)` - タスク完了

### 対話・ツール使用
- `LogUserInteraction(userMessage, claudeResponse)` - ユーザー対話
- `LogToolUsage(toolName, purpose, parameters)` - ツール使用
- `LogCodeGeneration(fileName, language, description)` - コード生成

### デバッグ
- `LogDebugInfo(context, details)` - デバッグ情報

## カスタマイズ

### ログレベルの変更
```go
// ロガー作成時にレベルを指定
vibeLogger := logger.New(logger.WARN) // WARNレベル以上のみ
```

### ログファイル名のカスタマイズ
```go
// カスタムファイル名
logFileName := "my_custom_log.txt"
logFilePath := filepath.Join(currentDir, logFileName)
```

### フォーマッターの変更
```go
// JSONフォーマッターを使用
jsonFormatter := formatter.NewJSONFormatter()
fileWriter.SetFormatter(jsonFormatter)
```

## 注意事項

- ログファイルは現在の作業ディレクトリに作成されます
- ファイル名には実行時の日時が含まれるため、複数実行しても上書きされません
- ログファイルは自動的に閉じられますが、明示的に`Close()`を呼ぶことを推奨します
- ログファイルが大きくなる場合は、ローテーション機能の使用を検討してください

## 高度な機能

### ログローテーション
大量のログを扱う場合は、ローテーション機能を使用できます：

```go
// サイズベースローテーション（10MB、最大5ファイル）
rotatingWriter, err := writer.NewRotatingFileWriter("claude_log.txt", 10*1024*1024, 5)

// 日次ローテーション
dailyWriter, err := writer.NewDailyRotatingFileWriter("claude_log")
```

### バッファリング
パフォーマンスを向上させるためのバッファリング：

```go
// 100エントリまでバッファリング
bufferedWriter, err := writer.NewBufferedFileWriter("claude_log.txt", 100)
```