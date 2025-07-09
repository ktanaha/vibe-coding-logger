package logger

import (
	"context"
	"time"
)

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

// Entry はログエントリの構造を表す
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
	
	// システム情報
	SystemInfo  map[string]interface{} `json:"system_info,omitempty"`
	RuntimeInfo map[string]interface{} `json:"runtime_info,omitempty"`
}

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

// OperationTracker は操作の追跡を行う
type OperationTracker struct {
	ID         string
	Operation  string
	StartTime  time.Time
	Input      map[string]interface{}
	Context    map[string]interface{}
	Logger     Logger
	parent     *OperationTracker
}

// Logger はロギングのインターフェースを定義する
type Logger interface {
	// 基本的なログ出力
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	// 操作の追跡
	StartOperation(operation string, input map[string]interface{}) *OperationTracker
	CompleteOperation(tracker *OperationTracker, output map[string]interface{})
	ErrorOperation(tracker *OperationTracker, err error, resolution string)

	// エラーハンドリング
	LogError(err error, context map[string]interface{}, retryable bool)
	LogRetry(operation string, attempt int, err error, nextRetryIn time.Duration)
	LogRecovery(operation string, originalErr error, recoveryAction string)

	// コンテキスト付きロガー
	WithContext(ctx context.Context) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithTag(tag string) Logger
	WithTags(tags []string) Logger

	// トレーシング
	WithTraceID(traceID string) Logger
	WithSpanID(spanID string) Logger
	WithParentID(parentID string) Logger

	// 設定
	SetLevel(level LogLevel)
	GetLevel() LogLevel
	AddWriter(writer Writer)
	SetFormatter(formatter Formatter)
	
	// システム情報設定
	EnableSystemInfo(enabled bool)
	EnableRuntimeInfo(enabled bool)
	IsSystemInfoEnabled() bool
	IsRuntimeInfoEnabled() bool
}

// Writer はログの出力先を定義する
type Writer interface {
	Write(entry *Entry) error
	Close() error
}

// Formatter はログのフォーマットを定義する
type Formatter interface {
	Format(entry *Entry) ([]byte, error)
}

// Field はログフィールドを表す
type Field struct {
	Key   string
	Value interface{}
}

// NewField は新しいフィールドを作成する
func NewField(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// String は文字列フィールドを作成する
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Int は整数フィールドを作成する
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Duration は期間フィールドを作成する
func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value}
}

// Error はエラーフィールドを作成する
func Error(key string, err error) Field {
	return Field{Key: key, Value: err.Error()}
}

// Any は任意の値のフィールドを作成する
func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}