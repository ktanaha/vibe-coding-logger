// Package internal は内部パッケージ間で共有される型定義を提供します。
package internal

import "time"

// Entry は内部用のログエントリ定義
type Entry struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Level       LogLevel               `json:"level"`
	Action      ActionType             `json:"action"`
	Operation   string                 `json:"operation"`
	Input       map[string]interface{} `json:"input,omitempty"`
	Output      map[string]interface{} `json:"output,omitempty"`
	Error       *ErrorInfo             `json:"error,omitempty"`
	Duration    time.Duration          `json:"duration,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	TraceID     string                 `json:"trace_id,omitempty"`
	SpanID      string                 `json:"span_id,omitempty"`
	ParentID    string                 `json:"parent_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	SystemInfo  map[string]interface{} `json:"system_info,omitempty"`
	RuntimeInfo map[string]interface{} `json:"runtime_info,omitempty"`
}

// LogLevel はログのレベルを表す
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String はLogLevelを文字列に変換する
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ActionType はアクションの種類を表す
type ActionType string

const (
	ActionStart    ActionType = "START"
	ActionComplete ActionType = "COMPLETE"
	ActionError    ActionType = "ERROR"
	ActionRetry    ActionType = "RETRY"
	ActionSkip     ActionType = "SKIP"
)

// ErrorInfo はエラー情報を表す
type ErrorInfo struct {
	Message    string                 `json:"message"`
	Type       string                 `json:"type"`
	Code       string                 `json:"code,omitempty"`
	Stack      string                 `json:"stack,omitempty"`
	Retryable  bool                   `json:"retryable"`
	Resolution string                 `json:"resolution,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`
}

// Formatter はログのフォーマットを定義する
type Formatter interface {
	Format(entry *Entry) ([]byte, error)
}

// Writer はログの出力先を定義する
type Writer interface {
	Write(entry *Entry) error
	Close() error
}
