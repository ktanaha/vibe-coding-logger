package logger

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// vibeLogger はLoggerインターフェースの実装
type vibeLogger struct {
	level     LogLevel
	writers   []Writer
	formatter Formatter
	fields    map[string]interface{}
	tags      []string
	traceID   string
	spanID    string
	parentID  string
	context   context.Context
	mu        sync.RWMutex

	// システム情報設定
	includeSystemInfo   bool
	includeRuntimeInfo  bool
	systemInfoCollector *SystemInfoCollector
}

// New は新しいloggerを作成する
func New(level LogLevel) Logger {
	return &vibeLogger{
		level:               level,
		writers:             make([]Writer, 0),
		fields:              make(map[string]interface{}),
		tags:                make([]string, 0),
		includeSystemInfo:   true,  // デフォルトで有効
		includeRuntimeInfo:  false, // パフォーマンスを考慮してデフォルトで無効
		systemInfoCollector: NewSystemInfoCollector(),
	}
}

// Default はデフォルトのロガーを作成する
func Default() Logger {
	logger := New(INFO)
	logger.AddWriter(NewConsoleWriter())
	logger.SetFormatter(NewTextFormatter())
	return logger
}

// Debug はデバッグレベルのログを出力する
func (l *vibeLogger) Debug(msg string, fields ...Field) {
	if l.level <= DEBUG {
		l.log(DEBUG, msg, fields...)
	}
}

// Info は情報レベルのログを出力する
func (l *vibeLogger) Info(msg string, fields ...Field) {
	if l.level <= INFO {
		l.log(INFO, msg, fields...)
	}
}

// Warn は警告レベルのログを出力する
func (l *vibeLogger) Warn(msg string, fields ...Field) {
	if l.level <= WARN {
		l.log(WARN, msg, fields...)
	}
}

// Error はエラーレベルのログを出力する
func (l *vibeLogger) Error(msg string, fields ...Field) {
	if l.level <= ERROR {
		l.log(ERROR, msg, fields...)
	}
}

// Fatal は致命的エラーレベルのログを出力する
func (l *vibeLogger) Fatal(msg string, fields ...Field) {
	if l.level <= FATAL {
		l.log(FATAL, msg, fields...)
	}
}

// log は実際のログ出力を行う
func (l *vibeLogger) log(level LogLevel, msg string, fields ...Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	entry := &Entry{
		ID:        uuid.New().String(),
		Timestamp: time.Now(),
		Level:     level,
		Operation: msg,
		Context:   l.copyFields(),
		Tags:      l.copyTags(),
		TraceID:   l.traceID,
		SpanID:    l.spanID,
		ParentID:  l.parentID,
		Metadata:  make(map[string]interface{}),
	}

	// システム情報の追加
	if l.includeSystemInfo {
		entry.SystemInfo = l.systemInfoCollector.GetCompactSystemInfo()
	}

	// ランタイム情報の追加
	if l.includeRuntimeInfo {
		entry.RuntimeInfo = l.systemInfoCollector.GetRuntimeStats()
	}

	// フィールドを追加
	for _, field := range fields {
		entry.Context[field.Key] = field.Value
	}

	// 呼び出し元の情報を追加
	if pc, file, line, ok := runtime.Caller(2); ok {
		entry.Metadata["caller"] = fmt.Sprintf("%s:%d", file, line)
		if fn := runtime.FuncForPC(pc); fn != nil {
			entry.Metadata["function"] = fn.Name()
		}
	}

	// 書き込み
	l.writeEntry(entry)
}

// StartOperation は操作の開始を記録する
func (l *vibeLogger) StartOperation(operation string, input map[string]interface{}) *OperationTracker {
	tracker := &OperationTracker{
		ID:        uuid.New().String(),
		Operation: operation,
		StartTime: time.Now(),
		Input:     input,
		Context:   l.copyFields(),
		Logger:    l,
	}

	l.log(INFO, operation,
		String("action", string(ActionStart)),
		String("operation_id", tracker.ID),
		Any("input", input))

	return tracker
}

// CompleteOperation は操作の完了を記録する
func (l *vibeLogger) CompleteOperation(tracker *OperationTracker, output map[string]interface{}) {
	duration := time.Since(tracker.StartTime)

	l.log(INFO, tracker.Operation,
		String("action", string(ActionComplete)),
		String("operation_id", tracker.ID),
		Any("input", tracker.Input),
		Any("output", output),
		Duration("duration", duration))
}

// ErrorOperation は操作のエラーを記録する
func (l *vibeLogger) ErrorOperation(tracker *OperationTracker, err error, resolution string) {
	duration := time.Since(tracker.StartTime)

	errorInfo := &ErrorInfo{
		Message:    err.Error(),
		Type:       fmt.Sprintf("%T", err),
		Retryable:  false,
		Resolution: resolution,
	}

	l.log(ERROR, tracker.Operation,
		String("action", string(ActionError)),
		String("operation_id", tracker.ID),
		Any("input", tracker.Input),
		Any("error", errorInfo),
		Duration("duration", duration))
}

// LogError はエラーをログに記録する
func (l *vibeLogger) LogError(err error, context map[string]interface{}, retryable bool) {
	errorInfo := &ErrorInfo{
		Message:   err.Error(),
		Type:      fmt.Sprintf("%T", err),
		Retryable: retryable,
		Context:   context,
	}

	// スタックトレースを取得
	errorInfo.Stack = l.getStackTrace()

	l.log(ERROR, "error_occurred",
		Any("error", errorInfo),
		Any("context", context))
}

// LogRetry はリトライの記録を行う
func (l *vibeLogger) LogRetry(operation string, attempt int, err error, nextRetryIn time.Duration) {
	l.log(WARN, operation,
		String("action", string(ActionRetry)),
		Int("attempt", attempt),
		Error("error", err),
		Duration("next_retry_in", nextRetryIn))
}

// LogRecovery はリカバリーの記録を行う
func (l *vibeLogger) LogRecovery(operation string, originalErr error, recoveryAction string) {
	l.log(INFO, operation,
		String("action", "RECOVERY"),
		Error("original_error", originalErr),
		String("recovery_action", recoveryAction))
}

// WithContext はコンテキスト付きロガーを返す
func (l *vibeLogger) WithContext(ctx context.Context) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := l.clone()
	newLogger.context = ctx
	return newLogger
}

// WithField はフィールド付きロガーを返す
func (l *vibeLogger) WithField(key string, value interface{}) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := l.clone()
	newLogger.fields[key] = value
	return newLogger
}

// WithFields は複数フィールド付きロガーを返す
func (l *vibeLogger) WithFields(fields map[string]interface{}) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := l.clone()
	for k, v := range fields {
		newLogger.fields[k] = v
	}
	return newLogger
}

// WithTag はタグ付きロガーを返す
func (l *vibeLogger) WithTag(tag string) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := l.clone()
	newLogger.tags = append(newLogger.tags, tag)
	return newLogger
}

// WithTags は複数タグ付きロガーを返す
func (l *vibeLogger) WithTags(tags []string) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := l.clone()
	newLogger.tags = append(newLogger.tags, tags...)
	return newLogger
}

// WithTraceID はトレースID付きロガーを返す
func (l *vibeLogger) WithTraceID(traceID string) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := l.clone()
	newLogger.traceID = traceID
	return newLogger
}

// WithSpanID はスパンID付きロガーを返す
func (l *vibeLogger) WithSpanID(spanID string) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := l.clone()
	newLogger.spanID = spanID
	return newLogger
}

// WithParentID は親ID付きロガーを返す
func (l *vibeLogger) WithParentID(parentID string) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := l.clone()
	newLogger.parentID = parentID
	return newLogger
}

// SetLevel はログレベルを設定する
func (l *vibeLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel は現在のログレベルを取得する
func (l *vibeLogger) GetLevel() LogLevel {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

// AddWriter はライターを追加する
func (l *vibeLogger) AddWriter(writer Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writers = append(l.writers, writer)
}

// SetFormatter はフォーマッターを設定する
func (l *vibeLogger) SetFormatter(formatter Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.formatter = formatter
}

// EnableSystemInfo はシステム情報の記録を有効/無効にする
func (l *vibeLogger) EnableSystemInfo(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.includeSystemInfo = enabled
}

// EnableRuntimeInfo はランタイム情報の記録を有効/無効にする
func (l *vibeLogger) EnableRuntimeInfo(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.includeRuntimeInfo = enabled
}

// IsSystemInfoEnabled はシステム情報の記録が有効かどうかを返す
func (l *vibeLogger) IsSystemInfoEnabled() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.includeSystemInfo
}

// IsRuntimeInfoEnabled はランタイム情報の記録が有効かどうかを返す
func (l *vibeLogger) IsRuntimeInfoEnabled() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.includeRuntimeInfo
}

// clone はロガーのクローンを作成する
func (l *vibeLogger) clone() *vibeLogger {
	newLogger := &vibeLogger{
		level:               l.level,
		writers:             l.writers,
		formatter:           l.formatter,
		fields:              make(map[string]interface{}),
		tags:                make([]string, 0),
		traceID:             l.traceID,
		spanID:              l.spanID,
		parentID:            l.parentID,
		context:             l.context,
		includeSystemInfo:   l.includeSystemInfo,
		includeRuntimeInfo:  l.includeRuntimeInfo,
		systemInfoCollector: l.systemInfoCollector,
	}

	// フィールドをコピー
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// タグをコピー
	newLogger.tags = append(newLogger.tags, l.tags...)

	return newLogger
}

// copyFields はフィールドのコピーを作成する
func (l *vibeLogger) copyFields() map[string]interface{} {
	fields := make(map[string]interface{})
	for k, v := range l.fields {
		fields[k] = v
	}
	return fields
}

// copyTags はタグのコピーを作成する
func (l *vibeLogger) copyTags() []string {
	tags := make([]string, len(l.tags))
	copy(tags, l.tags)
	return tags
}

// writeEntry はエントリを書き込む
func (l *vibeLogger) writeEntry(entry *Entry) {
	for _, writer := range l.writers {
		if err := writer.Write(entry); err != nil {
			// ライターでエラーが発生した場合の処理
			fmt.Printf("Logger write error: %v\n", err)
		}
	}
}

// getStackTrace はスタックトレースを取得する
func (l *vibeLogger) getStackTrace() string {
	var stack []string
	for i := 2; i < 10; i++ {
		if pc, file, line, ok := runtime.Caller(i); ok {
			if fn := runtime.FuncForPC(pc); fn != nil {
				stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
			}
		} else {
			break
		}
	}
	return strings.Join(stack, "\n")
}
