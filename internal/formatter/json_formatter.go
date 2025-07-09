package formatter

import (
	"encoding/json"
	"time"
	"vibe-coding-logger/pkg/logger"
)

// エイリアスを定義
type Entry = logger.Entry
type LogLevel = logger.LogLevel
type ActionType = logger.ActionType
type ErrorInfo = logger.ErrorInfo

// 定数のエイリアス
const (
	DEBUG = logger.DEBUG
	INFO  = logger.INFO
	WARN  = logger.WARN
	ERROR = logger.ERROR
	FATAL = logger.FATAL
)

// JSONFormatter はJSON形式でログを出力する
type JSONFormatter struct {
	TimestampFormat string
	PrettyPrint     bool
}

// NewJSONFormatter は新しいJSONフォーマッターを作成する
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     false,
	}
}

// NewPrettyJSONFormatter は整形されたJSONフォーマッターを作成する
func NewPrettyJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     true,
	}
}

// Format はエントリをJSON形式にフォーマットする
func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	// タイムスタンプを文字列に変換
	data := map[string]interface{}{
		"id":        entry.ID,
		"timestamp": entry.Timestamp.Format(f.TimestampFormat),
		"level":     entry.Level.String(),
		"action":    string(entry.Action),
		"operation": entry.Operation,
	}

	// 追加フィールドを設定
	if entry.Input != nil && len(entry.Input) > 0 {
		data["input"] = entry.Input
	}

	if entry.Output != nil && len(entry.Output) > 0 {
		data["output"] = entry.Output
	}

	if entry.Error != nil {
		data["error"] = entry.Error
	}

	if entry.Duration > 0 {
		data["duration"] = entry.Duration.String()
		data["duration_ms"] = entry.Duration.Milliseconds()
	}

	if entry.Context != nil && len(entry.Context) > 0 {
		data["context"] = entry.Context
	}

	if entry.Tags != nil && len(entry.Tags) > 0 {
		data["tags"] = entry.Tags
	}

	if entry.TraceID != "" {
		data["trace_id"] = entry.TraceID
	}

	if entry.SpanID != "" {
		data["span_id"] = entry.SpanID
	}

	if entry.ParentID != "" {
		data["parent_id"] = entry.ParentID
	}

	if entry.Metadata != nil && len(entry.Metadata) > 0 {
		data["metadata"] = entry.Metadata
	}

	// JSON形式でエンコード
	if f.PrettyPrint {
		return json.MarshalIndent(data, "", "  ")
	}

	return json.Marshal(data)
}

// VibeJSONFormatter はバイブコーディング専用のJSONフォーマッター
type VibeJSONFormatter struct {
	*JSONFormatter
	IncludeSessionInfo bool
	IncludeMetrics     bool
}

// NewVibeJSONFormatter は新しいバイブJSONフォーマッターを作成する
func NewVibeJSONFormatter() *VibeJSONFormatter {
	return &VibeJSONFormatter{
		JSONFormatter:      NewJSONFormatter(),
		IncludeSessionInfo: true,
		IncludeMetrics:     true,
	}
}

// Format はエントリをバイブコーディング用JSON形式にフォーマットする
func (f *VibeJSONFormatter) Format(entry *Entry) ([]byte, error) {
	// 基本のJSONフォーマットを取得
	baseData := make(map[string]interface{})
	baseBytes, err := f.JSONFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(baseBytes, &baseData); err != nil {
		return nil, err
	}

	// バイブコーディング専用フィールドを追加
	if f.IncludeSessionInfo {
		if sessionID, ok := getFromContext(entry.Context, "session_id"); ok {
			baseData["session_id"] = sessionID
		}
		if problemDomain, ok := getFromContext(entry.Context, "problem_domain"); ok {
			baseData["problem_domain"] = problemDomain
		}
		if programmingStep, ok := getFromContext(entry.Context, "programming_step"); ok {
			baseData["programming_step"] = programmingStep
		}
	}

	if f.IncludeMetrics {
		// メトリクス情報を追加
		metrics := make(map[string]interface{})
		
		if entry.Duration > 0 {
			metrics["performance"] = map[string]interface{}{
				"duration_ms":     entry.Duration.Milliseconds(),
				"duration_string": entry.Duration.String(),
			}
		}

		if len(metrics) > 0 {
			baseData["metrics"] = metrics
		}
	}

	// 再エンコード
	if f.PrettyPrint {
		return json.MarshalIndent(baseData, "", "  ")
	}

	return json.Marshal(baseData)
}

// getFromContext はコンテキストから値を取得する
func getFromContext(context map[string]interface{}, key string) (interface{}, bool) {
	if context == nil {
		return nil, false
	}
	value, exists := context[key]
	return value, exists
}

// CompactJSONFormatter はコンパクトなJSON形式でログを出力する
type CompactJSONFormatter struct {
	TimestampFormat string
	ExcludeFields   []string
}

// NewCompactJSONFormatter は新しいコンパクトJSONフォーマッターを作成する
func NewCompactJSONFormatter() *CompactJSONFormatter {
	return &CompactJSONFormatter{
		TimestampFormat: time.RFC3339,
		ExcludeFields:   []string{},
	}
}

// Format はエントリをコンパクトなJSON形式にフォーマットする
func (f *CompactJSONFormatter) Format(entry *Entry) ([]byte, error) {
	data := map[string]interface{}{
		"ts":  entry.Timestamp.Format(f.TimestampFormat),
		"lvl": entry.Level.String(),
		"op":  entry.Operation,
	}

	// 必要なフィールドのみを含める
	if entry.Error != nil {
		data["err"] = entry.Error.Message
	}

	if entry.Duration > 0 {
		data["dur"] = entry.Duration.Milliseconds()
	}

	if entry.TraceID != "" {
		data["tid"] = entry.TraceID
	}

	// 除外フィールドをチェック
	for _, excludeField := range f.ExcludeFields {
		delete(data, excludeField)
	}

	return json.Marshal(data)
}

// StructuredJSONFormatter は構造化されたJSON形式でログを出力する
type StructuredJSONFormatter struct {
	TimestampFormat string
	NestedFields    bool
}

// NewStructuredJSONFormatter は新しい構造化JSONフォーマッターを作成する
func NewStructuredJSONFormatter() *StructuredJSONFormatter {
	return &StructuredJSONFormatter{
		TimestampFormat: time.RFC3339,
		NestedFields:    true,
	}
}

// Format はエントリを構造化JSON形式にフォーマットする
func (f *StructuredJSONFormatter) Format(entry *Entry) ([]byte, error) {
	data := map[string]interface{}{
		"log": map[string]interface{}{
			"id":        entry.ID,
			"timestamp": entry.Timestamp.Format(f.TimestampFormat),
			"level":     entry.Level.String(),
			"message":   entry.Operation,
		},
	}

	// アクション情報
	if entry.Action != "" {
		data["action"] = map[string]interface{}{
			"type": string(entry.Action),
		}
	}

	// 入出力情報
	if entry.Input != nil && len(entry.Input) > 0 {
		data["input"] = entry.Input
	}

	if entry.Output != nil && len(entry.Output) > 0 {
		data["output"] = entry.Output
	}

	// エラー情報
	if entry.Error != nil {
		data["error"] = map[string]interface{}{
			"message":    entry.Error.Message,
			"type":       entry.Error.Type,
			"retryable":  entry.Error.Retryable,
			"code":       entry.Error.Code,
			"resolution": entry.Error.Resolution,
		}
		if entry.Error.Stack != "" {
			data["error"].(map[string]interface{})["stack"] = entry.Error.Stack
		}
	}

	// パフォーマンス情報
	if entry.Duration > 0 {
		data["performance"] = map[string]interface{}{
			"duration_ms":     entry.Duration.Milliseconds(),
			"duration_string": entry.Duration.String(),
		}
	}

	// トレーシング情報
	if entry.TraceID != "" || entry.SpanID != "" || entry.ParentID != "" {
		tracing := make(map[string]interface{})
		if entry.TraceID != "" {
			tracing["trace_id"] = entry.TraceID
		}
		if entry.SpanID != "" {
			tracing["span_id"] = entry.SpanID
		}
		if entry.ParentID != "" {
			tracing["parent_id"] = entry.ParentID
		}
		data["tracing"] = tracing
	}

	// メタデータ
	if entry.Metadata != nil && len(entry.Metadata) > 0 {
		data["metadata"] = entry.Metadata
	}

	// コンテキスト
	if entry.Context != nil && len(entry.Context) > 0 {
		data["context"] = entry.Context
	}

	// タグ
	if entry.Tags != nil && len(entry.Tags) > 0 {
		data["tags"] = entry.Tags
	}

	return json.Marshal(data)
}