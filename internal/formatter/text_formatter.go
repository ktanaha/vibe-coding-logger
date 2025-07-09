package formatter

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"vibe-coding-logger/internal"
)

// TextFormatter はテキスト形式でログを出力する
type TextFormatter struct {
	TimestampFormat string
	ColorEnabled    bool
	FullTimestamp   bool
	ShowCaller      bool
	ShowDuration    bool
	ShowTraceID     bool
	FieldSeparator  string
}

// NewTextFormatter は新しいテキストフォーマッターを作成する
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		TimestampFormat: time.RFC3339,
		ColorEnabled:    true,
		FullTimestamp:   true,
		ShowCaller:      true,
		ShowDuration:    true,
		ShowTraceID:     true,
		FieldSeparator:  " ",
	}
}

// Format はエントリをテキスト形式にフォーマットする
func (f *TextFormatter) Format(entry *internal.Entry) ([]byte, error) {
	var parts []string

	// タイムスタンプ
	if f.FullTimestamp {
		parts = append(parts, entry.Timestamp.Format(f.TimestampFormat))
	} else {
		parts = append(parts, entry.Timestamp.Format("15:04:05"))
	}

	// ログレベル
	levelStr := f.formatLevel(entry.Level)
	parts = append(parts, levelStr)

	// トレースID
	if f.ShowTraceID && entry.TraceID != "" {
		parts = append(parts, fmt.Sprintf("[%s]", entry.TraceID))
	}

	// メッセージ
	parts = append(parts, entry.Operation)

	// アクション
	if entry.Action != "" {
		parts = append(parts, fmt.Sprintf("action=%s", entry.Action))
	}

	// 期間
	if f.ShowDuration && entry.Duration > 0 {
		parts = append(parts, fmt.Sprintf("duration=%s", entry.Duration))
	}

	// コンテキストフィールド
	if len(entry.Context) > 0 {
		contextStr := f.formatFields(entry.Context)
		if contextStr != "" {
			parts = append(parts, contextStr)
		}
	}

	// エラー
	if entry.Error != nil {
		errorStr := fmt.Sprintf("error=%q", entry.Error.Message)
		if entry.Error.Code != "" {
			errorStr += fmt.Sprintf(" code=%s", entry.Error.Code)
		}
		parts = append(parts, errorStr)
	}

	// 入力
	if len(entry.Input) > 0 {
		parts = append(parts, fmt.Sprintf("input=%s", f.formatMapCompact(entry.Input)))
	}

	// 出力
	if len(entry.Output) > 0 {
		parts = append(parts, fmt.Sprintf("output=%s", f.formatMapCompact(entry.Output)))
	}

	// タグ
	if len(entry.Tags) > 0 {
		parts = append(parts, fmt.Sprintf("tags=%s", strings.Join(entry.Tags, ",")))
	}

	// 呼び出し元
	if f.ShowCaller && len(entry.Metadata) > 0 {
		if caller, ok := entry.Metadata["caller"]; ok {
			parts = append(parts, fmt.Sprintf("caller=%s", caller))
		}
	}

	result := strings.Join(parts, f.FieldSeparator)
	return []byte(result + "\n"), nil
}

// formatLevel はログレベルを色付きでフォーマットする
func (f *TextFormatter) formatLevel(level internal.LogLevel) string {
	levelStr := level.String()
	
	if !f.ColorEnabled {
		return fmt.Sprintf("[%s]", levelStr)
	}

	switch level {
	case internal.DEBUG:
		return fmt.Sprintf("\033[36m[%s]\033[0m", levelStr) // Cyan
	case internal.INFO:
		return fmt.Sprintf("\033[32m[%s]\033[0m", levelStr) // Green
	case internal.WARN:
		return fmt.Sprintf("\033[33m[%s]\033[0m", levelStr) // Yellow
	case internal.ERROR:
		return fmt.Sprintf("\033[31m[%s]\033[0m", levelStr) // Red
	case internal.FATAL:
		return fmt.Sprintf("\033[35m[%s]\033[0m", levelStr) // Magenta
	default:
		return fmt.Sprintf("[%s]", levelStr)
	}
}

// formatFields はフィールドをキー=値の形式でフォーマットする
func (f *TextFormatter) formatFields(fields map[string]interface{}) string {
	if len(fields) == 0 {
		return ""
	}

	var parts []string
	keys := make([]string, 0, len(fields))
	
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := fields[k]
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}

	return strings.Join(parts, f.FieldSeparator)
}

// formatMapCompact はマップをコンパクト形式でフォーマットする
func (f *TextFormatter) formatMapCompact(m map[string]interface{}) string {
	if len(m) == 0 {
		return "{}"
	}

	var parts []string
	for k, v := range m {
		parts = append(parts, fmt.Sprintf("%s:%v", k, v))
	}

	return "{" + strings.Join(parts, ",") + "}"
}

// ConsoleFormatter はコンソール出力用のフォーマッター
type ConsoleFormatter struct {
	*TextFormatter
	UseEmoji bool
}

// NewConsoleFormatter は新しいコンソールフォーマッターを作成する
func NewConsoleFormatter() *ConsoleFormatter {
	return &ConsoleFormatter{
		TextFormatter: NewTextFormatter(),
		UseEmoji:      true,
	}
}

// Format はエントリをコンソール向けにフォーマットする
func (f *ConsoleFormatter) Format(entry *internal.Entry) ([]byte, error) {
	var parts []string

	// タイムスタンプ（短縮版）
	parts = append(parts, entry.Timestamp.Format("15:04:05"))

	// ログレベル（絵文字付き）
	levelStr := f.formatLevelWithEmoji(entry.Level)
	parts = append(parts, levelStr)

	// メッセージ
	parts = append(parts, entry.Operation)

	// 重要な情報のみ表示
	if entry.Error != nil {
		parts = append(parts, fmt.Sprintf("❌ %s", entry.Error.Message))
	}

	if entry.Duration > 0 {
		parts = append(parts, fmt.Sprintf("⏱️  %s", entry.Duration))
	}

	// バイブコーディング情報
	if len(entry.Context) > 0 {
		if sessionID, ok := entry.Context["session_id"]; ok {
			parts = append(parts, fmt.Sprintf("🔧 %s", sessionID))
		}
		if step, ok := entry.Context["programming_step"]; ok {
			parts = append(parts, fmt.Sprintf("📝 %s", step))
		}
	}

	result := strings.Join(parts, " ")
	return []byte(result + "\n"), nil
}

// formatLevelWithEmoji はログレベルを絵文字付きでフォーマットする
func (f *ConsoleFormatter) formatLevelWithEmoji(level internal.LogLevel) string {
	if !f.UseEmoji {
		return f.formatLevel(level)
	}

	switch level {
	case internal.DEBUG:
		return "🔍 DEBUG"
	case internal.INFO:
		return "ℹ️  INFO"
	case internal.WARN:
		return "⚠️  WARN"
	case internal.ERROR:
		return "❌ ERROR"
	case internal.FATAL:
		return "💀 FATAL"
	default:
		return fmt.Sprintf("❓ %s", level.String())
	}
}

// VibeTextFormatter はバイブコーディング専用のテキストフォーマッター
type VibeTextFormatter struct {
	*TextFormatter
	ShowSessionInfo bool
	ShowStepInfo    bool
	UseIcons        bool
}

// NewVibeTextFormatter は新しいバイブテキストフォーマッターを作成する
func NewVibeTextFormatter() *VibeTextFormatter {
	return &VibeTextFormatter{
		TextFormatter:   NewTextFormatter(),
		ShowSessionInfo: true,
		ShowStepInfo:    true,
		UseIcons:        true,
	}
}

// Format はエントリをバイブコーディング用テキスト形式にフォーマットする
func (f *VibeTextFormatter) Format(entry *internal.Entry) ([]byte, error) {
	var parts []string

	// タイムスタンプ
	parts = append(parts, entry.Timestamp.Format("15:04:05"))

	// ログレベル
	levelStr := f.formatLevel(entry.Level)
	parts = append(parts, levelStr)

	// セッション情報
	if f.ShowSessionInfo && len(entry.Context) > 0 {
		if sessionID, ok := entry.Context["session_id"]; ok {
			icon := "🔧"
			if f.UseIcons {
				parts = append(parts, fmt.Sprintf("%s[%s]", icon, sessionID))
			} else {
				parts = append(parts, fmt.Sprintf("[%s]", sessionID))
			}
		}
	}

	// ステップ情報
	if f.ShowStepInfo && len(entry.Context) > 0 {
		if step, ok := entry.Context["programming_step"]; ok {
			icon := f.getStepIcon(fmt.Sprintf("%v", step))
			if f.UseIcons {
				parts = append(parts, fmt.Sprintf("%s %s", icon, step))
			} else {
				parts = append(parts, fmt.Sprintf("<%s>", step))
			}
		}
	}

	// メッセージ
	parts = append(parts, entry.Operation)

	// 特別な情報の表示
	if entry.Error != nil {
		icon := "❌"
		if f.UseIcons {
			parts = append(parts, fmt.Sprintf("%s %s", icon, entry.Error.Message))
		} else {
			parts = append(parts, fmt.Sprintf("ERROR: %s", entry.Error.Message))
		}
	}

	if entry.Duration > 0 {
		icon := "⏱️"
		if f.UseIcons {
			parts = append(parts, fmt.Sprintf("%s %s", icon, entry.Duration))
		} else {
			parts = append(parts, fmt.Sprintf("(%s)", entry.Duration))
		}
	}

	// その他のコンテキスト情報
	if len(entry.Context) > 0 {
		filteredContext := make(map[string]interface{})
		for k, v := range entry.Context {
			if k != "session_id" && k != "programming_step" {
				filteredContext[k] = v
			}
		}
		if len(filteredContext) > 0 {
			contextStr := f.formatFields(filteredContext)
			if contextStr != "" {
				parts = append(parts, contextStr)
			}
		}
	}

	result := strings.Join(parts, " ")
	return []byte(result + "\n"), nil
}

// getStepIcon はプログラミングステップに応じたアイコンを取得する
func (f *VibeTextFormatter) getStepIcon(step string) string {
	switch strings.ToLower(step) {
	case "thinking", "思考":
		return "🤔"
	case "coding", "コーディング":
		return "💻"
	case "testing", "テスト":
		return "🧪"
	case "debugging", "デバッグ":
		return "🐛"
	case "refactoring", "リファクタリング":
		return "♻️"
	case "learning", "学習":
		return "📚"
	case "planning", "計画":
		return "📋"
	case "implementing", "実装":
		return "🔨"
	case "reviewing", "レビュー":
		return "👀"
	case "documenting", "文書化":
		return "📝"
	default:
		return "📍"
	}
}

// CompactTextFormatter はコンパクトなテキスト形式でログを出力する
type CompactTextFormatter struct {
	TimestampFormat string
	ColorEnabled    bool
}

// NewCompactTextFormatter は新しいコンパクトテキストフォーマッターを作成する
func NewCompactTextFormatter() *CompactTextFormatter {
	return &CompactTextFormatter{
		TimestampFormat: "15:04:05",
		ColorEnabled:    true,
	}
}

// Format はエントリをコンパクトテキスト形式にフォーマットする
func (f *CompactTextFormatter) Format(entry *internal.Entry) ([]byte, error) {
	var parts []string

	// タイムスタンプ
	parts = append(parts, entry.Timestamp.Format(f.TimestampFormat))

	// ログレベル
	levelStr := entry.Level.String()
	if f.ColorEnabled {
		switch entry.Level {
		case internal.ERROR:
			levelStr = fmt.Sprintf("\033[31m%s\033[0m", levelStr)
		case internal.WARN:
			levelStr = fmt.Sprintf("\033[33m%s\033[0m", levelStr)
		case internal.INFO:
			levelStr = fmt.Sprintf("\033[32m%s\033[0m", levelStr)
		}
	}
	parts = append(parts, levelStr)

	// メッセージ
	parts = append(parts, entry.Operation)

	// エラーがある場合のみ表示
	if entry.Error != nil {
		parts = append(parts, fmt.Sprintf("err=%q", entry.Error.Message))
	}

	result := strings.Join(parts, " ")
	return []byte(result + "\n"), nil
}