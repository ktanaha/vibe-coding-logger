package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	logger "vibe-coding-logger/pkg/logger"
)

// ClaudeCodeLogger はClaude Code専用のロガー設定
type ClaudeCodeLogger struct {
	logger     logger.Logger
	fileWriter logger.Writer
	logFile    string
}

// NewClaudeCodeLogger は新しいClaude Code用ロガーを作成する
func NewClaudeCodeLogger() (*ClaudeCodeLogger, error) {
	// 現在の作業ディレクトリを取得
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("作業ディレクトリの取得に失敗しました: %v", err)
	}

	// ログファイル名を現在の日時で生成
	timestamp := time.Now().Format("20060102_150405")
	logFileName := fmt.Sprintf("claude_code_log_%s.txt", timestamp)
	logFilePath := filepath.Join(currentDir, logFileName)

	// ファイルライターを作成
	fileWriter, err := logger.NewFileWriter(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルライターの作成に失敗しました: %v", err)
	}

	// テキストフォーマッターを設定
	textFormatter := logger.NewTextFormatter()

	// ファイルライターにフォーマッターを設定
	if fwImpl, ok := fileWriter.(interface{ SetFormatter(logger.Formatter) }); ok {
		fwImpl.SetFormatter(textFormatter)
	}

	// ロガーを作成
	vibeLogger := logger.New(logger.DEBUG)
	vibeLogger.AddWriter(fileWriter)
	vibeLogger.SetFormatter(textFormatter)

	// コンソール出力も有効にする場合は以下をアンコメント
	// vibeLogger.AddWriter(logger.NewConsoleWriter())

	ccLogger := &ClaudeCodeLogger{
		logger:     vibeLogger,
		fileWriter: fileWriter,
		logFile:    logFilePath,
	}

	// 初期ログを出力
	ccLogger.LogClaudeCodeStart()

	return ccLogger, nil
}

// LogClaudeCodeStart はClaude Code開始をログに記録する
func (c *ClaudeCodeLogger) LogClaudeCodeStart() {
	c.logger.Info("Claude Code セッション開始",
		logger.String("log_file", c.logFile),
		logger.String("start_time", time.Now().Format("2006-01-02 15:04:05")),
	)
}

// LogCommand はコマンド実行をログに記録する
func (c *ClaudeCodeLogger) LogCommand(command string, args []string) {
	c.logger.Info("コマンド実行",
		logger.String("command", command),
		logger.Any("arguments", args),
		logger.String("timestamp", time.Now().Format("2006-01-02 15:04:05")),
	)
}

// LogFileOperation はファイル操作をログに記録する
func (c *ClaudeCodeLogger) LogFileOperation(operation, filePath string) {
	c.logger.Info("ファイル操作",
		logger.String("operation", operation),
		logger.String("file_path", filePath),
		logger.String("timestamp", time.Now().Format("2006-01-02 15:04:05")),
	)
}

// LogError はエラーをログに記録する
func (c *ClaudeCodeLogger) LogError(operation string, err error) {
	c.logger.Error("エラー発生",
		logger.String("operation", operation),
		logger.Error("error", err),
		logger.String("timestamp", time.Now().Format("2006-01-02 15:04:05")),
	)
}

// LogTaskStart はタスク開始をログに記録する
func (c *ClaudeCodeLogger) LogTaskStart(taskName, description string) {
	c.logger.Info("タスク開始",
		logger.String("task_name", taskName),
		logger.String("description", description),
		logger.String("timestamp", time.Now().Format("2006-01-02 15:04:05")),
	)
}

// LogTaskComplete はタスク完了をログに記録する
func (c *ClaudeCodeLogger) LogTaskComplete(taskName string, duration time.Duration) {
	c.logger.Info("タスク完了",
		logger.String("task_name", taskName),
		logger.Duration("duration", duration),
		logger.String("timestamp", time.Now().Format("2006-01-02 15:04:05")),
	)
}

// LogUserInteraction はユーザーとの対話をログに記録する
func (c *ClaudeCodeLogger) LogUserInteraction(userMessage, claudeResponse string) {
	c.logger.Info("ユーザー対話",
		logger.String("user_message", userMessage),
		logger.String("claude_response", claudeResponse),
		logger.String("timestamp", time.Now().Format("2006-01-02 15:04:05")),
	)
}

// LogToolUsage はツール使用をログに記録する
func (c *ClaudeCodeLogger) LogToolUsage(toolName, purpose string, parameters map[string]interface{}) {
	c.logger.Info("ツール使用",
		logger.String("tool_name", toolName),
		logger.String("purpose", purpose),
		logger.Any("parameters", parameters),
		logger.String("timestamp", time.Now().Format("2006-01-02 15:04:05")),
	)
}

// LogCodeGeneration はコード生成をログに記録する
func (c *ClaudeCodeLogger) LogCodeGeneration(fileName, language, description string) {
	c.logger.Info("コード生成",
		logger.String("file_name", fileName),
		logger.String("language", language),
		logger.String("description", description),
		logger.String("timestamp", time.Now().Format("2006-01-02 15:04:05")),
	)
}

// LogDebugInfo はデバッグ情報をログに記録する
func (c *ClaudeCodeLogger) LogDebugInfo(context string, details map[string]interface{}) {
	c.logger.Debug("デバッグ情報",
		logger.String("context", context),
		logger.Any("details", details),
		logger.String("timestamp", time.Now().Format("2006-01-02 15:04:05")),
	)
}

// GetLogFilePath はログファイルのパスを返す
func (c *ClaudeCodeLogger) GetLogFilePath() string {
	return c.logFile
}

// Close はロガーを閉じる
func (c *ClaudeCodeLogger) Close() error {
	c.logger.Info("Claude Code セッション終了",
		logger.String("end_time", time.Now().Format("2006-01-02 15:04:05")),
	)
	
	if c.fileWriter != nil {
		return c.fileWriter.Close()
	}
	return nil
}

// 使用例を示すmain関数
func main() {
	// Claude Codeロガーを初期化
	ccLogger, err := NewClaudeCodeLogger()
	if err != nil {
		log.Fatalf("ロガーの初期化に失敗しました: %v", err)
	}
	defer ccLogger.Close()

	fmt.Printf("ログファイルが作成されました: %s\n", ccLogger.GetLogFilePath())

	// 使用例
	ccLogger.LogTaskStart("ファイル作成", "新しいGoファイルを作成します")
	
	ccLogger.LogFileOperation("CREATE", "example.go")
	
	ccLogger.LogCommand("go", []string{"build", "."})
	
	ccLogger.LogUserInteraction(
		"Goファイルを作成してください",
		"新しいGoファイルを作成しました。ビルドも正常に完了しました。",
	)
	
	ccLogger.LogToolUsage("Write", "ファイル作成", map[string]interface{}{
		"file_path": "example.go",
		"content_length": 150,
	})
	
	ccLogger.LogCodeGeneration("example.go", "Go", "サンプルのGo言語プログラム")
	
	ccLogger.LogTaskComplete("ファイル作成", time.Second*2)
	
	ccLogger.LogDebugInfo("プロセス状況", map[string]interface{}{
		"memory_usage": "15MB",
		"cpu_usage": "2.5%",
		"active_goroutines": 5,
	})

	fmt.Println("ログ出力のサンプルが完了しました。")
}