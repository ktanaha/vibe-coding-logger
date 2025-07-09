package logger

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// ErrorHandler はエラーハンドリングに特化した機能を提供する
type ErrorHandler struct {
	logger Logger
}

// NewErrorHandler は新しいエラーハンドラーを作成する
func NewErrorHandler(logger Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// HandleError はエラーを処理し、適切なログを出力する
func (eh *ErrorHandler) HandleError(err error, context map[string]interface{}, options ...ErrorOption) {
	if err == nil {
		return
	}

	errorInfo := &ErrorInfo{
		Message:   err.Error(),
		Type:      fmt.Sprintf("%T", err),
		Retryable: false,
		Context:   context,
		Stack:     eh.getStackTrace(2),
	}

	// オプションを適用
	for _, option := range options {
		option(errorInfo)
	}

	eh.logger.Error("error_handled",
		String("error_type", errorInfo.Type),
		String("error_message", errorInfo.Message),
		String("error_code", errorInfo.Code),
		String("resolution", errorInfo.Resolution),
		String("stack_trace", errorInfo.Stack),
		Any("context", errorInfo.Context),
		Any("retryable", errorInfo.Retryable))
}

// HandlePanic はパニックを処理する
func (eh *ErrorHandler) HandlePanic(recovered interface{}, context map[string]interface{}) {
	if recovered == nil {
		return
	}

	errorInfo := &ErrorInfo{
		Message:   fmt.Sprintf("Panic recovered: %v", recovered),
		Type:      "panic",
		Retryable: false,
		Context:   context,
		Stack:     eh.getStackTrace(2),
	}

	eh.logger.Fatal("panic_recovered",
		String("panic_value", fmt.Sprintf("%v", recovered)),
		String("stack_trace", errorInfo.Stack),
		Any("context", errorInfo.Context))
}

// HandleRetryableError はリトライ可能なエラーを処理する
func (eh *ErrorHandler) HandleRetryableError(err error, attempt int, maxAttempts int, nextRetryIn time.Duration, context map[string]interface{}) {
	errorInfo := &ErrorInfo{
		Message:   err.Error(),
		Type:      fmt.Sprintf("%T", err),
		Retryable: true,
		Context:   context,
	}

	eh.logger.Warn("retryable_error",
		String("error_type", errorInfo.Type),
		String("error_message", errorInfo.Message),
		Int("attempt", attempt),
		Int("max_attempts", maxAttempts),
		Duration("next_retry_in", nextRetryIn),
		Any("context", errorInfo.Context))
}

// HandleFatalError は致命的なエラーを処理する
func (eh *ErrorHandler) HandleFatalError(err error, context map[string]interface{}, shutdownFunc func()) {
	if err == nil {
		return
	}

	errorInfo := &ErrorInfo{
		Message:   err.Error(),
		Type:      fmt.Sprintf("%T", err),
		Retryable: false,
		Context:   context,
		Stack:     eh.getStackTrace(2),
	}

	eh.logger.Fatal("fatal_error",
		String("error_type", errorInfo.Type),
		String("error_message", errorInfo.Message),
		String("stack_trace", errorInfo.Stack),
		Any("context", errorInfo.Context))

	if shutdownFunc != nil {
		shutdownFunc()
	}
}

// HandleRecovery はリカバリーを処理する
func (eh *ErrorHandler) HandleRecovery(originalErr error, recoveryAction string, recoveryResult interface{}, context map[string]interface{}) {
	eh.logger.Info("recovery_executed",
		String("original_error", originalErr.Error()),
		String("recovery_action", recoveryAction),
		Any("recovery_result", recoveryResult),
		Any("context", context))
}

// getStackTrace はスタックトレースを取得する
func (eh *ErrorHandler) getStackTrace(skip int) string {
	var stack []string
	for i := skip; i < skip+10; i++ {
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

// ErrorOption はエラー情報を設定するためのオプション
type ErrorOption func(*ErrorInfo)

// WithErrorCode はエラーコードを設定する
func WithErrorCode(code string) ErrorOption {
	return func(ei *ErrorInfo) {
		ei.Code = code
	}
}

// WithRetryable はリトライ可能性を設定する
func WithRetryable(retryable bool) ErrorOption {
	return func(ei *ErrorInfo) {
		ei.Retryable = retryable
	}
}

// WithResolution は解決策を設定する
func WithResolution(resolution string) ErrorOption {
	return func(ei *ErrorInfo) {
		ei.Resolution = resolution
	}
}

// WithContext は追加のコンテキストを設定する
func WithContext(key string, value interface{}) ErrorOption {
	return func(ei *ErrorInfo) {
		if ei.Context == nil {
			ei.Context = make(map[string]interface{})
		}
		ei.Context[key] = value
	}
}

// RetryHandler はリトライ処理に特化したハンドラー
type RetryHandler struct {
	logger       Logger
	errorHandler *ErrorHandler
}

// NewRetryHandler は新しいリトライハンドラーを作成する
func NewRetryHandler(logger Logger) *RetryHandler {
	return &RetryHandler{
		logger:       logger,
		errorHandler: NewErrorHandler(logger),
	}
}

// ExecuteWithRetry は指定した関数をリトライ付きで実行する
func (rh *RetryHandler) ExecuteWithRetry(operation string, fn func() error, maxAttempts int, backoff time.Duration, context map[string]interface{}) error {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		rh.logger.Info(operation,
			String("action", "ATTEMPT"),
			Int("attempt", attempt),
			Int("max_attempts", maxAttempts),
			Any("context", context))

		err := fn()
		if err == nil {
			// 成功
			rh.logger.Info(operation,
				String("action", "SUCCESS"),
				Int("attempt", attempt),
				Any("context", context))
			return nil
		}

		lastErr = err

		if attempt < maxAttempts {
			// リトライ
			nextRetryIn := backoff * time.Duration(attempt)
			rh.errorHandler.HandleRetryableError(err, attempt, maxAttempts, nextRetryIn, context)
			time.Sleep(nextRetryIn)
		} else {
			// 最大試行回数に達した
			rh.errorHandler.HandleError(err, context,
				WithErrorCode("MAX_RETRIES_EXCEEDED"),
				WithResolution("手動での対応が必要です"))
		}
	}

	return lastErr
}

// ExecuteWithCircuitBreaker はサーキットブレーカー付きで実行する
func (rh *RetryHandler) ExecuteWithCircuitBreaker(operation string, fn func() error, failureThreshold int, context map[string]interface{}) error {
	// シンプルなサーキットブレーカーの実装
	// 実際の実装では、より洗練された状態管理が必要

	rh.logger.Info(operation,
		String("action", "CIRCUIT_BREAKER_EXECUTE"),
		Int("failure_threshold", failureThreshold),
		Any("context", context))

	err := fn()
	if err != nil {
		rh.errorHandler.HandleError(err, context,
			WithErrorCode("CIRCUIT_BREAKER_FAILURE"),
			WithResolution("サーキットブレーカーが作動しました"))
	}

	return err
}

// RecoveryHandler はリカバリー処理に特化したハンドラー
type RecoveryHandler struct {
	logger       Logger
	errorHandler *ErrorHandler
}

// NewRecoveryHandler は新しいリカバリーハンドラーを作成する
func NewRecoveryHandler(logger Logger) *RecoveryHandler {
	return &RecoveryHandler{
		logger:       logger,
		errorHandler: NewErrorHandler(logger),
	}
}

// ExecuteWithRecovery はリカバリー付きで実行する
func (rh *RecoveryHandler) ExecuteWithRecovery(operation string, fn func() error, recoveryFn func(error) error, context map[string]interface{}) error {
	rh.logger.Info(operation,
		String("action", "EXECUTE_WITH_RECOVERY"),
		Any("context", context))

	err := fn()
	if err != nil {
		rh.logger.Warn(operation,
			String("action", "RECOVERY_NEEDED"),
			String("error", err.Error()),
			Any("context", context))

		if recoveryFn != nil {
			recoveryErr := recoveryFn(err)
			if recoveryErr != nil {
				rh.errorHandler.HandleError(recoveryErr, context,
					WithErrorCode("RECOVERY_FAILED"),
					WithResolution("リカバリーが失敗しました"))
				return recoveryErr
			}

			rh.errorHandler.HandleRecovery(err, "custom_recovery", "success", context)
		}
	}

	return err
}

// ExecuteWithPanicRecovery はパニックリカバリー付きで実行する
func (rh *RecoveryHandler) ExecuteWithPanicRecovery(operation string, fn func(), context map[string]interface{}) (recovered interface{}) {
	defer func() {
		if r := recover(); r != nil {
			recovered = r
			rh.errorHandler.HandlePanic(r, context)
		}
	}()

	rh.logger.Info(operation,
		String("action", "EXECUTE_WITH_PANIC_RECOVERY"),
		Any("context", context))

	fn()
	return nil
}

// VibeErrorHandler はバイブコーディング専用のエラーハンドラー
type VibeErrorHandler struct {
	*ErrorHandler
	sessionID     string
	problemDomain string
}

// NewVibeErrorHandler は新しいバイブエラーハンドラーを作成する
func NewVibeErrorHandler(logger Logger, sessionID, problemDomain string) *VibeErrorHandler {
	return &VibeErrorHandler{
		ErrorHandler:  NewErrorHandler(logger),
		sessionID:     sessionID,
		problemDomain: problemDomain,
	}
}

// HandleCodingError はコーディングエラーを処理する
func (veh *VibeErrorHandler) HandleCodingError(err error, codeFile string, lineNumber int, codeSnippet string, resolution string) {
	context := map[string]interface{}{
		"session_id":     veh.sessionID,
		"problem_domain": veh.problemDomain,
		"code_file":      codeFile,
		"line_number":    lineNumber,
		"code_snippet":   codeSnippet,
	}

	veh.HandleError(err, context,
		WithErrorCode("CODING_ERROR"),
		WithResolution(resolution))
}

// HandleTestError はテストエラーを処理する
func (veh *VibeErrorHandler) HandleTestError(err error, testName string, testOutput string, expectedVsActual string) {
	context := map[string]interface{}{
		"session_id":         veh.sessionID,
		"problem_domain":     veh.problemDomain,
		"test_name":          testName,
		"test_output":        testOutput,
		"expected_vs_actual": expectedVsActual,
	}

	veh.HandleError(err, context,
		WithErrorCode("TEST_ERROR"),
		WithResolution("テストコードまたは実装コードの修正が必要です"))
}

// HandleBuildError はビルドエラーを処理する
func (veh *VibeErrorHandler) HandleBuildError(err error, buildOutput string, dependencies []string) {
	context := map[string]interface{}{
		"session_id":     veh.sessionID,
		"problem_domain": veh.problemDomain,
		"build_output":   buildOutput,
		"dependencies":   dependencies,
	}

	veh.HandleError(err, context,
		WithErrorCode("BUILD_ERROR"),
		WithResolution("依存関係の確認またはビルド設定の修正が必要です"))
}

// HandleLogicError はロジックエラーを処理する
func (veh *VibeErrorHandler) HandleLogicError(err error, expectation string, actual string, debugInfo map[string]interface{}) {
	context := map[string]interface{}{
		"session_id":     veh.sessionID,
		"problem_domain": veh.problemDomain,
		"expectation":    expectation,
		"actual":         actual,
		"debug_info":     debugInfo,
	}

	veh.HandleError(err, context,
		WithErrorCode("LOGIC_ERROR"),
		WithResolution("アルゴリズムまたはロジックの見直しが必要です"))
}
